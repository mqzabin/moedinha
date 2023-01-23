package currency

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// decimalPlaces defines how many decimal digits should be used.
	currDecimalDigits = 10
	// currMaxIntegerDigits ...
	currMaxIntegerDigits = (natNumberOfInts * natMaxDigitsPerInt) - currDecimalDigits
	//
	currDecimalPointerSymbol = '.'
	//
	currNegativeSymbol = '-'
)

var (
	zeroFiller = strings.Repeat("0", natNumberOfInts*natMaxDigitsPerInt)

	currencyRegexp = regexp.MustCompile(fmt.Sprintf(
		`^-?\d{1,%d}(\%c\d{0,%d})?$`,
		currMaxIntegerDigits,
		currDecimalPointerSymbol,
		currDecimalDigits,
	))
)

type Currency struct {
	n   nat
	neg bool
}

func FromDecimalString(v string) (Currency, error) {
	if !currencyRegexp.MatchString(v) {
		return Currency{}, errors.New("invalid currency format")
	}

	var neg bool
	if v[0] == currNegativeSymbol {
		neg = true
		v = v[1:]
	}

	// pointIndex (right to left)
	pointIndex := strings.IndexRune(v, currDecimalPointerSymbol)

	var decimalDigits int
	if pointIndex >= 0 {
		decimalDigits = len(v) - (pointIndex + 1)
		v = v[:pointIndex] + v[pointIndex+1:]
	}

	n1Pow := currDecimalDigits - decimalDigits
	if n1Pow > currDecimalDigits {
		panic("decimal digits overflow")
	}

	leftZeros := (natMaxDigitsPerInt * natNumberOfInts) - (len(v) + n1Pow)

	var builder strings.Builder
	builder.Grow(natMaxDigitsPerInt * natNumberOfInts)
	builder.WriteString(zeroFiller[:leftZeros])
	builder.WriteString(v)
	builder.WriteString(zeroFiller[:n1Pow])

	n, err := newNatFromString(builder.String())
	if err != nil {
		return Currency{}, fmt.Errorf("creating natural number from string")
	}

	return Currency{
		n:   n,
		neg: neg,
	}, nil
}

func (c Currency) String() string {

	str := c.n.string()
	str = fmt.Sprintf("%s.%s", str[:currMaxIntegerDigits], str[currMaxIntegerDigits:])

	//return str

	var leftZerosToRemove int
	for _, ch := range str {
		if ch != '0' {
			break
		}

		leftZerosToRemove++
	}

	var rightZerosToRemove int
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] != '0' {
			break
		}

		rightZerosToRemove++
	}

	if str[len(str)-rightZerosToRemove-1] == '.' {
		rightZerosToRemove++
	}

	if str[leftZerosToRemove] == '.' {
		leftZerosToRemove--
	}

	return str[leftZerosToRemove : len(str)-rightZerosToRemove]
}

func (c Currency) IsZero() bool {
	return c.n.isZero()
}

func (c Currency) Equal(v Currency) bool {
	// -0 should be equal to +0.
	if c.n.isZero() && c.n.isZero() {
		return true
	}

	return c.neg == v.neg && c.n.equal(v.n)
}

func (c Currency) GreaterThan(v Currency) bool {
	if c.n.isZero() && c.n.isZero() {
		return false
	}

	if c.neg == v.neg {
		return c.n.greaterThan(v.n) && !c.neg
	}

	if !c.neg {
		return true
	}

	return false
}

func (c Currency) Add(v Currency) Currency {
	if c.neg == v.neg {
		return Currency{
			n:   c.n.add(v.n),
			neg: c.neg,
		}
	}

	if c.neg {
		return v.Sub(c)
	}

	return c.Sub(v)
}

// TODO: Draft
func (c Currency) Sub(v Currency) Currency {
	if c.Equal(v) {
		return Currency{}
	}

	// different signals
	if c.neg != v.neg {
		// c - (-v) = c + v
		if v.neg {
			return Currency{
				n: c.n.add(v.n),
			}
		}

		// -c - v = - (c+v)
		return Currency{
			n:   c.n.add(v.n),
			neg: true,
		}
	}

	// for now on, equal sign

	// both negative
	// -c - (-v) = v - c
	if c.neg {
		if c.n.greaterThan(v.n) {
			// v - c = -(c-v)
			return Currency{
				n:   c.n.sub(v.n),
				neg: true,
			}
		}

		return Currency{
			n: v.n.sub(c.n),
		}
	}

	// both positive
	// c - v

	if v.n.greaterThan(c.n) {
		// c - v = -(v - c)
		return Currency{
			n:   v.n.sub(c.n),
			neg: true,
		}
	}

	return Currency{
		n: c.n.sub(v.n),
	}
}
