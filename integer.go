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
	abs natural
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
		abs: n,
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

	natString := t.abs.string()

	copy(intString[1:], natString[:])

	return intString
}

// add sum two integers.
func (t integer) add(v integer) integer {
	// "(+t)+(+v) = t+v" or "(-t)+(-v) = -(t+v)"
	if t.neg == v.neg {
		return integer{
			abs: t.abs.add(v.abs),
			neg: t.neg,
		}
	}

	// For now on, signs are different.

	// C is negative.
	// v - t
	if t.neg {
		return v.sub(integer{
			abs: t.abs,
			neg: false,
		})
	}

	// V is negative.
	// t - v

	return t.sub(integer{
		abs: v.abs,
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
				abs: t.abs.add(v.abs),
				neg: false,
			}
		}

		// -c - v = - (c+v)
		return integer{
			abs: t.abs.add(v.abs),
			neg: true,
		}
	}

	// for now on, equal sign

	// both negative numbers
	// -t - (-v) = v - t
	if t.neg {
		// negative result.
		if t.abs.greaterThan(v.abs) {
			// v - t = -(t-v)
			return integer{
				abs: t.abs.sub(v.abs),
				neg: true,
			}
		}

		// positive result
		return integer{
			abs: v.abs.sub(t.abs),
			neg: false,
		}
	}

	// both positive
	// c - v

	// negative result
	if v.abs.greaterThan(t.abs) {
		// t - v = -(v - t)
		return integer{
			abs: v.abs.sub(t.abs),
			neg: true,
		}
	}

	// positive result
	return integer{
		abs: t.abs.sub(v.abs),
		neg: false,
	}
}

// mul multiplies two integer numbers.
// The first return is the result, and the second return is the overflow
// of the operation, if any, as a natural number.
func (t integer) mul(v integer) (integer, natural) {
	natResult, natOverflow := t.abs.mul(v.abs)

	return integer{
		abs: natResult,
		neg: t.neg != v.neg,
	}, natOverflow
}

func (t integer) div(v integer) (integer, natural) {
	q, r := t.abs.div(v.abs)

	return integer{
		abs: q,
		neg: v.neg != t.neg,
	}, r
}

func (t integer) divByInt(v int64) (integer, uint64) {
	negV := v < 0
	if negV {
		v *= -1
	}

	q, r := t.abs.divByUint(uint64(v))

	return integer{
		abs: q,
		neg: negV != t.neg,
	}, r
}

func (t integer) isZero() bool {
	return t.abs.isZero()
}

func (t integer) equal(v integer) bool {
	// -0 should be equal to +0.
	if t.abs.isZero() && v.abs.isZero() {
		return true
	}

	return t.neg == v.neg && t.abs.equal(v.abs)
}

func (t integer) greaterThan(v integer) bool {
	if t.isZero() && v.isZero() {
		return false
	}

	// equal signal
	if neg := t.neg; neg == v.neg {
		if neg {
			return t.abs.lessThan(v.abs)
		}

		return t.abs.greaterThan(v.abs)
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
			return t.abs.lessThan(v.abs)
		}

		return t.abs.greaterThan(v.abs)
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
			return t.abs.greaterThan(v.abs)
		}

		return t.abs.lessThan(v.abs)
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
			return t.abs.greaterThan(v.abs)
		}

		return t.abs.lessThan(v.abs)
	}

	if t.neg {
		return true
	}

	return false
}

func (t integer) leftShiftUint(toShift int) (integer, natural) {
	shifted, overflow := t.abs.leftShiftUint(toShift)

	return integer{
		abs: shifted,
		neg: t.neg,
	}, overflow
}

func (t integer) rightShiftUint(toShift int) (integer, natural) {
	shifted, overflow := t.abs.rightShiftUint(toShift)

	return integer{
		abs: shifted,
		neg: t.neg,
	}, overflow
}
