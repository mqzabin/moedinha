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
	t integer
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

	return Currency{t: intValue}, nil
}

func (c Currency) String() string {
	if c.t.isZero() {
		return "0"
	}

	intString := c.t.string()

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

	if c.t.neg {
		currString[leftZerosToRemove] = integerNegativeSymbol
		leftZerosToRemove--
	}

	return string(currString[leftZerosToRemove+1 : len(currString)-rightZerosToRemove])
}

func (c Currency) IsZero() bool {
	return c.t.isZero()
}

func (c Currency) Equal(v Currency) bool {
	return c.t.equal(v.t)
}

func (c Currency) GreaterThan(v Currency) bool {
	return c.t.greaterThan(v.t)
}

func (c Currency) GreaterThanOrEqual(v Currency) bool {
	return c.t.greaterThanOrEqual(v.t)
}

func (c Currency) LessThan(v Currency) bool {
	return c.t.lessThan(v.t)
}

func (c Currency) LessThanOrEqual(v Currency) bool {
	return c.t.lessThanOrEqual(v.t)
}

func (c Currency) Add(v Currency) Currency {
	return Currency{c.t.add(v.t)}
}

func (c Currency) Sub(v Currency) Currency {
	return Currency{c.t.sub(v.t)}
}

func (c Currency) Mul(v Currency) Currency {
	intResult, natOverflow := c.t.mul(v.t)

	// Since integers and naturals represents numbers with currencyDecimalDigits decimal
	// digits, the result represents a number with 2*currencyDecimalDigits decimal digits.
	// There's a need to truncate the first currencyDecimalDigits from the natural number.
	intResult.n, _ = intResult.n.rightShiftUint(uintsReservedToDecimal)

	// Getting the overflow part that should be summed to result.
	natOverflow, addToResult := natOverflow.rightShiftUint(uintsReservedToDecimal)

	intResult.n = intResult.n.add(addToResult)

	if !natOverflow.isZero() {
		panic(fmt.Sprintf("multiplication overflow: %s * %s", c.String(), v.String()))
	}

	return Currency{
		t: intResult,
	}
}

//func (c Currency) Div(v Currency) Currency {
//	if c.Equal(v) {
//		return currencyOne
//	}
//
//	if c.GreaterThan(v) {
//		return c.divWithInitialEstimate(initialEstimateChebyshev, v)
//	}
//
//	// c is less than v
//
//	if c.LessThan(v) && !c.Equal(currencyOne) {
//		return currencyOne.divWithInitialEstimate(initialEstimate, v.divWithInitialEstimate(initialEstimate, c))
//	}
//
//	if c.Equal(currencyOne) {
//
//	}
//
//	return c.divWithInitialEstimate(initialEstimateChebyshev, v)
//}

func (c Currency) divWithInitialEstimate(initialEstimate func(Currency) Currency, v Currency) Currency {
	if v.IsZero() {
		panic("dividing by zero")
	}

	nDigits := v.t.n.digits()

	neededShift := nDigits - currencyDecimalDigits

	cShift, _ := c.t.n.rightShiftDigit(neededShift)
	vShift, _ := v.t.n.rightShiftDigit(neededShift)

	for vShift.lessThan(currencyOneHalf.t.n) {
		vShift, _ = vShift.mulByUint64(2)
		cShift, _ = cShift.mulByUint64(2)
	}

	shiftedNumerator := Currency{t: integer{
		n:   cShift,
		neg: false,
	}}

	shiftedDenominator := Currency{t: integer{
		n:   vShift,
		neg: false,
	}}

	reciprocal := initialEstimate(shiftedDenominator)

	for {
		mul := shiftedDenominator.Mul(reciprocal)
		if mul.Equal(currencyOne) {
			//fmt.Println(reciprocal.String())
			//fmt.Println("*")
			//fmt.Println(shiftedNumerator.String())
			//fmt.Println("=")
			//fmt.Println(reciprocal.Mul(shiftedNumerator))
			break
		}

		reciprocal = reciprocal.Mul(currencyTwo.Sub(mul))
	}

	result := shiftedNumerator.Mul(reciprocal)
	result.t.neg = c.t.neg != v.t.neg

	return result
}

func initialEstimateOne(_ Currency) Currency {
	return currencyOne
}

func initialEstimateChebyshev(v Currency) Currency {
	// (48/17) - Denominator*(32/17)
	return currency48Over17.Sub(v.Mul(currency32Over17))
}
