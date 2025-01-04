package moedinha

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// currencyDecimalDigits defines how many decimal digits should be used.
	currencyDecimalDigits = uintsReservedToDecimal * maxDigitsPerUint
	// currencyMaxIntegerDigits the amount of digits before the decimal pointer.
	currencyMaxIntegerDigits = (numberOfUints * maxDigitsPerUint) - currencyDecimalDigits
	// currencyDecimalSeparatorSymbol the separator used for decimal digits.
	currencyDecimalSeparatorSymbol = '.'
	// currencyMaxLen stores the maximum length of a currency string.
	// +1 for possible decimal separator.
	currencyMaxLen = integerMaxLen + 1
)

var (
	ErrInvalidFormat = errors.New("invalid format")

	currencyRegexp = regexp.MustCompile(fmt.Sprintf(
		`^-?\d{1,%d}(\%c\d{0,%d})?$`,
		currencyMaxIntegerDigits,
		currencyDecimalSeparatorSymbol,
		currencyDecimalDigits,
	))
)

type Currency struct {
	mantissa integer
}

func NewFromString(str string) (Currency, error) {
	if !currencyRegexp.MatchString(str) {
		return Currency{}, fmt.Errorf(`validating currency: "%s": %w`, str, ErrInvalidFormat)
	}

	separatorIndex := strings.IndexRune(str, currencyDecimalSeparatorSymbol)

	strLen := len(str)

	decimalDigits := 0
	integerDigits := strLen
	if separatorIndex >= 0 {
		decimalDigits = strLen - (separatorIndex + 1)
		integerDigits -= decimalDigits + 1
	}

	// How much the copy should shift in integer number string. For example:
	// 0.1 is shifted to left by 17 digits if the support is for 18 digits.
	cpRightShift := currencyDecimalDigits - decimalDigits
	cpLeftShift := integerMaxLen - cpRightShift - (decimalDigits + integerDigits)

	intString := [integerMaxLen]byte{zeroRune}

	// Copy the integer part.
	copy(intString[cpLeftShift:cpLeftShift+integerDigits], str[:integerDigits])
	// Copying the decimal part, if any.
	if decimalDigits > 0 {
		copy(intString[cpLeftShift+integerDigits:cpLeftShift+integerDigits+decimalDigits], str[integerDigits+1:])
	}

	// Adding leading zeros.
	copy(intString[:cpLeftShift], zeroFiller[:])
	// Adding trailing zeros.
	copy(intString[currencyMaxLen-(cpRightShift+1):], zeroFiller[:])

	if cpLeftShift != 0 && intString[cpLeftShift] == integerNegativeSymbol {
		intString[0] = integerNegativeSymbol
		intString[cpLeftShift] = zeroRune
	}

	intValue, err := newIntegerFromString(intString)
	if err != nil {
		return Currency{}, fmt.Errorf("creating underlying integer: %w", err)
	}

	return Currency{mantissa: intValue}, nil
}

func (c Currency) String() string {
	if c.mantissa.isZero() {
		return "0"
	}

	intString := c.mantissa.string()

	currString := [currencyMaxLen]byte{}
	// Copy integer part, preserving 0 index to negative symbol.
	copy(currString[:currencyMaxIntegerDigits+1], intString[:currencyMaxIntegerDigits+1])
	// Setting the decimal separator.
	currString[currencyMaxIntegerDigits+1] = currencyDecimalSeparatorSymbol
	// Copying the decimal part.
	copy(currString[currencyMaxIntegerDigits+2:], intString[currencyMaxIntegerDigits+1:])

	var leftZerosToRemove int
	for i := 1; i < currencyMaxLen; i++ {
		if currString[i] != zeroRune {
			break
		}

		leftZerosToRemove++
	}

	var rightZerosToRemove int
	for i := currencyMaxLen - 1; i >= 0; i-- {
		if currString[i] != zeroRune {
			break
		}

		rightZerosToRemove++
	}

	// Removing decimal separator if it's an integer number, e.g: "1.0" turn into "1"
	if currString[currencyMaxLen-rightZerosToRemove-1] == currencyDecimalSeparatorSymbol {
		rightZerosToRemove++
	}

	// Preserving at least one 0 at left side of separator, e.g: ".1" turn into "0.1"
	if currString[leftZerosToRemove+1] == currencyDecimalSeparatorSymbol {
		leftZerosToRemove--
	}

	if c.mantissa.neg {
		currString[leftZerosToRemove] = integerNegativeSymbol
		leftZerosToRemove--
	}

	return string(currString[leftZerosToRemove+1 : len(currString)-rightZerosToRemove])
}

func (c Currency) IsZero() bool {
	return c.mantissa.isZero()
}

func (c Currency) Equal(v Currency) bool {
	return c.mantissa.equal(v.mantissa)
}

func (c Currency) GreaterThan(v Currency) bool {
	return c.mantissa.greaterThan(v.mantissa)
}

func (c Currency) GreaterThanOrEqual(v Currency) bool {
	return c.mantissa.greaterThanOrEqual(v.mantissa)
}

func (c Currency) LessThan(v Currency) bool {
	return c.mantissa.lessThan(v.mantissa)
}

func (c Currency) LessThanOrEqual(v Currency) bool {
	return c.mantissa.lessThanOrEqual(v.mantissa)
}

func (c Currency) Add(v Currency) Currency {
	return Currency{c.mantissa.add(v.mantissa)}
}

func (c Currency) Sub(v Currency) Currency {
	return Currency{c.mantissa.sub(v.mantissa)}
}

func (c Currency) Mul(v Currency) Currency {
	intResult, natOverflow := c.mantissa.mul(v.mantissa)

	// Since integers and naturals represents numbers with currencyDecimalDigits decimal
	// digits, the result represents a number with 2*currencyDecimalDigits decimal digits.
	// There's a need to truncate the first currencyDecimalDigits from the natural number.
	intResult.abs, _ = intResult.abs.rightShiftUint(uintsReservedToDecimal)

	// Getting the overflow part that should be summed to result.
	natOverflow, addToResult := natOverflow.rightShiftUint(uintsReservedToDecimal)

	intResult.abs = intResult.abs.add(addToResult)

	if !natOverflow.isZero() {
		panic(fmt.Sprintf("multiplication overflow: %s * %s", c.String(), v.String()))
	}

	return Currency{
		mantissa: intResult,
	}
}

func (c Currency) DivInt(v int64) Currency {
	intRes, _ := c.mantissa.divByInt(v)

	return Currency{
		mantissa: intRes,
	}
}

func (c Currency) Div(v Currency) Currency {
	var shiftScale int

	shouldRound := func(r natural) bool {
		r2, _ := r.mulByUint64(2)
		return r2.greaterThan(v.mantissa.abs)
	}

	// Checking if 'c' < 'v'. If so, we should shift 'c' to the left until
	// it's greater than v.
	if c.mantissa.abs.lessThan(v.mantissa.abs) {
		cUints := c.mantissa.abs.uintsInUse()
		vUints := v.mantissa.abs.uintsInUse()

		shiftScale = vUints - cUints
		if c.mantissa.abs[numberOfUints-cUints] < v.mantissa.abs[numberOfUints-vUints] {
			shiftScale++
		}
	}

	// Shifting and splitting dividend in its high and low parts.
	var dividend [highLow]natural
	switch shiftScale {
	case 0:
		dividend[low] = c.mantissa.abs
	default:
		dividend[low], dividend[high] = c.mantissa.abs.leftShiftUint(shiftScale)
	}

	var quotient, remainder [highLow]natural

	quotient[high], remainder[high] = dividend[high].div(v.mantissa.abs)

	quotient[low], remainder[low] = dividend[low].div(v.mantissa.abs)
	if shouldRound(remainder[low]) {
		quotient[low] = quotient[low].add(naturalOne)
	}

	// Adjusting the result decimal digits.
	quotient[low], _ = quotient[low].leftShiftUint(uintsReservedToDecimal)

	// Un-shifting result.
	if shiftScale > 0 {
		var underflow natural

		quotient[low], underflow = quotient[low].rightShiftUint(shiftScale)

		// Rounding if needed.
		if underflow[0] >= halfMaxValuePerUint+1 {
			quotient[low] = quotient[low].add(naturalOne)
		}
	}

	return Currency{
		mantissa: integer{
			abs: quotient[low],
			neg: c.mantissa.neg != v.mantissa.neg,
		},
	}
}

func (c Currency) DivRound(v Currency) Currency {
	v.mantissa.abs = v.mantissa.abs.sub(naturalOne)

	res := c.Div(v)

	res.mantissa.abs = res.mantissa.abs.add(naturalOne)

	return res
}
