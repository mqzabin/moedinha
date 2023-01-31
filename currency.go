package moedinha

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// currencyDecimalDigits defines how many decimal digits should be used.
	currencyDecimalDigits = 10
	// currencyMaxIntegerDigits the amount of digits before the decimal pointer.
	currencyMaxIntegerDigits = (natNumberOfInts * natMaxDigitsPerInt) - currencyDecimalDigits
	// currencyDecimalSeparatorSymbol the separator used for decimal digits.
	currencyDecimalSeparatorSymbol = '.'
	// currencyNegativeSymbol symbol used to represent a negative number as a string.
	currencyNegativeSymbol = '-'
	// maxCurrencyLen stores the maximum length of a currency string. +2 for possible decimal separator and
	// negative symbol
	maxCurrencyLen = natDigits + 2
)

var currencyRegexp = regexp.MustCompile(fmt.Sprintf(
	`^-?\d{1,%d}(\%c\d{0,%d})?$`,
	currencyMaxIntegerDigits,
	currencyDecimalSeparatorSymbol,
	currencyDecimalDigits,
))

type Currency struct {
	n   nat
	neg bool
}

func NewFromString(str string) (Currency, error) {
	if !currencyRegexp.MatchString(str) {
		return Currency{}, errors.New("invalid currency format")
	}

	var neg bool
	if str[0] == currencyNegativeSymbol {
		neg = true
		str = str[1:] // Ignore the negative symbol.
	}

	l := len(str)

	var natStr [natDigits]byte

	separatorIndex := strings.IndexRune(str, currencyDecimalSeparatorSymbol)

	decimalDigits := 0
	integerDigits := l
	if separatorIndex >= 0 {
		decimalDigits = l - (separatorIndex + 1)
		integerDigits = l - decimalDigits - 1 // -1 for the separator.
	}

	// How much the copy should shift in natural number string. For example:
	// 0.001 should be positioned as xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx0001xxxxxxx in natural number string.
	cpRightShift := currencyDecimalDigits - decimalDigits
	cpLeftShift := natDigits - (decimalDigits + integerDigits) - cpRightShift

	// Copying the integer part.
	copy(natStr[cpLeftShift:cpLeftShift+integerDigits], str[:integerDigits])
	// Copying the decimal part, if any.
	if decimalDigits > 0 {
		copy(natStr[cpLeftShift+integerDigits:cpLeftShift+integerDigits+decimalDigits], str[integerDigits+1:])
	}

	// Adding leading zeros.
	copy(natStr[:cpLeftShift], zeroFiller[:cpLeftShift])
	// Adding trailing zeros.
	copy(natStr[natDigits-cpRightShift:], zeroFiller[:cpRightShift])

	n, err := newNatFromString(natStr)
	if err != nil {
		return Currency{}, fmt.Errorf("creating natural number from string: %w", err)
	}

	return Currency{
		n:   n,
		neg: neg,
	}, nil
}

func (c Currency) String() string {
	if c.IsZero() {
		return "0"
	}

	natString := c.n.string()

	currString := [maxCurrencyLen]byte{zeroRune}
	// Copy integer part, preserving 0 index to negative symbol.
	copy(currString[1:currencyMaxIntegerDigits+1], natString[:currencyMaxIntegerDigits])
	// Setting the decimal separator.
	currString[currencyMaxIntegerDigits+1] = currencyDecimalSeparatorSymbol
	// Copying the decimal part.
	copy(currString[currencyMaxIntegerDigits+2:], natString[currencyMaxIntegerDigits:])

	var leftZerosToRemove int
	for i := 1; i < maxCurrencyLen; i++ {
		if currString[i] != zeroRune {
			break
		}

		leftZerosToRemove++
	}

	var rightZerosToRemove int
	for i := maxCurrencyLen - 1; i >= 0; i-- {
		if currString[i] != zeroRune {
			break
		}

		rightZerosToRemove++
	}

	// Removing decimal separator if it's an integer number, e.g. 1.0 turn into 1
	if currString[maxCurrencyLen-rightZerosToRemove-1] == currencyDecimalSeparatorSymbol {
		rightZerosToRemove++
	}

	// Preserving at least one 0 at left side of separator, e.g. .1 turn into 0.1
	if currString[leftZerosToRemove+1] == currencyDecimalSeparatorSymbol {
		leftZerosToRemove--
	}

	if c.neg {
		currString[leftZerosToRemove] = currencyNegativeSymbol
		leftZerosToRemove--
	}

	return string(currString[leftZerosToRemove+1 : len(currString)-rightZerosToRemove])
}

func (c Currency) IsZero() bool {
	return c.n.isZero()
}

func (c Currency) Equal(v Currency) bool {
	// -0 should be equal to +0.
	if c.n.isZero() && v.n.isZero() {
		return true
	}

	return c.neg == v.neg && c.n.equal(v.n)
}

func (c Currency) GreaterThan(v Currency) bool {
	if c.Equal(v) {
		return false
	}

	// Equal signal
	if neg := c.neg; neg == v.neg {
		if neg {
			return c.n.lessThan(v.n)
		}

		return c.n.greaterThan(v.n)
	}

	if c.neg {
		return false
	}

	return true
}

func (c Currency) LessThan(v Currency) bool {
	if c.Equal(v) {
		return false
	}
	// Equal signal
	if neg := c.neg; neg == v.neg {
		if neg {
			return c.n.greaterThan(v.n)
		}

		return c.n.lessThan(v.n)
	}

	if c.neg {
		return true
	}

	return false
}

func (c Currency) Add(v Currency) Currency {
	// "(+c)+(+v) = c+v" or "(-c)+(-v) = -(c+v)"
	if c.neg == v.neg {
		return Currency{
			n:   c.n.add(v.n),
			neg: c.neg,
		}
	}

	// For now on, signals are different.

	// C is negative.
	if c.neg {
		// (-c)+(+v) = (+v) - (+c)
		return v.Sub(Currency{
			n:   c.n,
			neg: false,
		})
	}

	// V is negative.

	// (+c)+(-v) = (+c) - (+v)
	return c.Sub(Currency{
		n:   v.n,
		neg: false,
	})
}

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
		// will be negative.
		if c.n.greaterThan(v.n) {
			// v - c = -(c-v)
			return Currency{
				n:   v.n.difference(c.n),
				neg: true,
			}
		}

		return Currency{
			n: c.n.difference(v.n),
		}
	}

	// both positive
	// c - v

	if v.n.greaterThan(c.n) {
		// c - v = -(v - c)
		return Currency{
			n:   c.n.difference(v.n),
			neg: true,
		}
	}

	return Currency{
		n: v.n.difference(c.n),
	}
}
