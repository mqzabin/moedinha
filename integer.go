package moedinha

import (
	"fmt"
)

const (
	// integerNegativeSymbol symbol used to represent a negative number as a string.
	integerNegativeSymbol = '-'
	// integerMaxLen is the maximum length that an integer string.
	// +1 to the possible negative symbol.
	integerMaxLen = naturalMaxLen + 1
)

type integer struct {
	n   natural
	neg bool
}

func newIntegerFromString(str [integerMaxLen]byte) (integer, error) {
	var neg bool
	if str[0] == integerNegativeSymbol {
		neg = true
	}

	var natStr [naturalMaxLen]byte

	copy(natStr[:], str[1:])

	n, err := newNatFromString(natStr)
	if err != nil {
		return integer{}, fmt.Errorf("creating underlyin natural number from string: %w", err)
	}

	return integer{
		n:   n,
		neg: neg,
	}, nil
}

func (t integer) string() [integerMaxLen]byte {
	intString := [integerMaxLen]byte{zeroRune}

	if t.isZero() {
		copy(intString[:], zeroFiller[:integerMaxLen])

		return intString
	}

	if t.neg {
		intString[0] = integerNegativeSymbol
	}

	natString := t.n.string()

	copy(intString[1:], natString[:])

	return intString
}

// add sum two integers.
func (t integer) add(v integer) integer {
	// "(+t)+(+v) = t+v" or "(-t)+(-v) = -(t+v)"
	if t.neg == v.neg {
		return integer{
			n:   t.n.add(v.n),
			neg: t.neg,
		}
	}

	// For now on, signs are different.

	// C is negative.
	// v - t
	if t.neg {
		return v.sub(integer{
			n:   t.n,
			neg: false,
		})
	}

	// V is negative.
	// t - v

	return t.sub(integer{
		n:   v.n,
		neg: false,
	})
}

// sub calculates the subtraction "t - v".
func (t integer) sub(v integer) integer {
	if t.equal(v) {
		return integer{}
	}

	// different signs
	if t.neg != v.neg {
		// t - (-v) = t + v
		if v.neg {
			return integer{
				n:   t.n.add(v.n),
				neg: false,
			}
		}

		// -c - v = - (c+v)
		return integer{
			n:   t.n.add(v.n),
			neg: true,
		}
	}

	// for now on, equal sign

	// both negative numbers
	// -t - (-v) = v - t
	if t.neg {
		// negative result.
		if t.n.greaterThan(v.n) {
			// v - t = -(t-v)
			return integer{
				n:   t.n.sub(v.n),
				neg: true,
			}
		}

		// positive result
		return integer{
			n:   v.n.sub(t.n),
			neg: false,
		}
	}

	// both positive
	// c - v

	// negative result
	if v.n.greaterThan(t.n) {
		// t - v = -(v - t)
		return integer{
			n:   v.n.sub(t.n),
			neg: true,
		}
	}

	// positive result
	return integer{
		n:   t.n.sub(v.n),
		neg: false,
	}
}

// mul multiplies two integer numbers.
// The first return is the result, and the second return is the overflow
// of the operation, if any, as a natural number.
func (t integer) mul(v integer) (integer, natural) {
	natResult, natOverflow := t.n.mul(v.n)

	if t.neg == v.neg {
		return integer{
			n:   natResult,
			neg: false,
		}, natOverflow
	}

	return integer{
		n:   natResult,
		neg: true,
	}, natOverflow
}

func (t integer) isZero() bool {
	return t.n.isZero()
}

func (t integer) equal(v integer) bool {
	// -0 should be equal to +0.
	if t.n.isZero() && v.n.isZero() {
		return true
	}

	return t.neg == v.neg && t.n.equal(v.n)
}

func (t integer) greaterThan(v integer) bool {
	if t.isZero() && v.isZero() {
		return false
	}

	// equal signal
	if neg := t.neg; neg == v.neg {
		if neg {
			return t.n.lessThan(v.n)
		}

		return t.n.greaterThan(v.n)
	}

	if t.neg {
		return false
	}

	return true
}

func (t integer) greaterThanOrEqual(v integer) bool {
	if t.equal(v) {
		return true
	}

	// equal signal
	if neg := t.neg; neg == v.neg {
		if neg {
			return t.n.lessThan(v.n)
		}

		return t.n.greaterThan(v.n)
	}

	if t.neg {
		return false
	}

	return true
}

func (t integer) lessThan(v integer) bool {
	if t.isZero() && v.isZero() {
		return false
	}

	// equal signs
	if neg := t.neg; neg == v.neg {
		if neg {
			return t.n.greaterThan(v.n)
		}

		return t.n.lessThan(v.n)
	}

	if t.neg {
		return true
	}

	return false
}

func (t integer) lessThanOrEqual(v integer) bool {
	if t.equal(v) {
		return true
	}

	// equal signs
	if neg := t.neg; neg == v.neg {
		if neg {
			return t.n.greaterThan(v.n)
		}

		return t.n.lessThan(v.n)
	}

	if t.neg {
		return true
	}

	return false
}
