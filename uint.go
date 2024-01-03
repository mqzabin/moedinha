package moedinha

import (
	"fmt"
)

const (
	// base is the base value used by the library.
	base = 10
	// maxValuePerUint is the greater 999-ish number under 63 bits.
	// Its binary representation is:
	// 101100110111110110011011111100000110100111100111111111111111
	maxValuePerUint = 999999999999999999
	// maxValuePerUintBitsUsed it's the number of bits used to represent maxValuePerUint.
	maxValuePerUintBitsUsed = 60
	// maxDigitsPerUint is the amount of maxValuePerUint digits.
	maxDigitsPerUint = 18
	// halfMaxValuePerUint is half of maxValuePerUint digits.
	halfMaxValuePerUint = 999999999
)

var pow10 = newPow10Map()

func newPow10Map() map[int]uint64 {
	powValue := uint64(1)

	pow10Map := make(map[int]uint64, 2*naturalMaxLen)
	pow10Map[0] = powValue

	for e := 1; e <= naturalMaxLen; e++ {
		powValue *= base

		pow10Map[e] = powValue
	}

	return pow10Map
}

// TODO: Implement some sort of binary search here.
func digitsOf(v uint64) int {
	if v == 0 {
		return 0
	}

	for i := maxDigitsPerUint; i >= 0; i-- {
		if v/pow10[i] != 0 {
			return i + 1
		}
	}

	panic("finding digits of a uint64")
}

func rightShift(v uint64, shift int) (uint64, uint64) {
	if shift > maxDigitsPerUint {
		panic("invalid shift passed to right shift")
	}

	if shift < 0 {
		return leftShift(v, -shift)
	}

	return v / pow10[shift], (v % pow10[shift]) * pow10[maxDigitsPerUint-shift]
}

func leftShift(v uint64, shift int) (uint64, uint64) {
	if shift > maxDigitsPerUint {
		panic("invalid shift passed to left shift")
	}

	if shift < 0 {
		return rightShift(v, -shift)
	}

	return (v % pow10[maxDigitsPerUint-shift]) * pow10[shift], v / pow10[maxDigitsPerUint-shift]
}

// rebalance truncates the src to maxValuePerUint, returns it at newSrc and adds the reminder to newDest.
func rebalance(src, dest uint64) (newSrc, newDest uint64) {
	dest += src / (maxValuePerUint + 1)
	src %= maxValuePerUint + 1

	return src, dest
}

// splitInHalf splits the given uint64 into two parts of 9 digits.
// The first return is the left part, and last is the right part.
func splitInHalf(a uint64) (uint64, uint64) {
	return a / (halfMaxValuePerUint + 1), a % (halfMaxValuePerUint + 1)
}

// multiplyUint multiply two uint64 and return the result and
// the overflow of the multiplication.
func multiplyUint(a, b uint64) (uint64, uint64) {
	// All those components are less than 32-bits values.
	aLeft, aRight := splitInHalf(a) // aLeft is the right bits of a, and al the right bits.
	bLeft, bRight := splitInHalf(b) // The same.

	// Each operation will overflow the 32-bits components and go to 64-bits components.
	// So right, middle, and left will potentially overflow the 32-bits boundary.
	// The idea here is that:
	// a = (aLeft.10^9 + aRight)
	// b = (bLeft.10^9 + bRight)
	// So a.b = aLeft.bLeft.10^(2*9) + (aLeft.bRight + aRight.bLeft).10^(9) + (aRight*bRight)
	right := aRight * bRight
	middle := aLeft*bRight + aRight*bLeft
	left := bLeft * aLeft

	// Now we should distribute the middle term to the respective side:
	middleLeft, middleRight := splitInHalf(middle)

	// middleLeft is added to left as the first digits.
	left += middleLeft
	// middleRight is added to right as the last digits.
	right += middleRight * (halfMaxValuePerUint + 1)

	// Rebalance right to left to avoid overflows.
	right, left = rebalance(right, left)

	return right, left
}

// atoi is a fork from strconv.Atoi with proper signature.
func atoi(s [maxDigitsPerUint]byte) (uint64, error) {
	var n uint64
	for _, ch := range s {
		ch -= zeroRune
		if ch > base-1 {
			return 0, fmt.Errorf("invalid syntax converting string to uint64: rune %c", ch)
		}
		n = n*base + uint64(ch)
	}

	return n, nil
}

// itoa is a fork from strconv.Itoa with proper signature.
func itoa(v uint64) [maxDigitsPerUint]byte {
	var res [maxDigitsPerUint]byte
	div := v

	for i := range res {
		digit := byte(div % base)
		res[maxDigitsPerUint-i-1] = zeroRune + digit

		div /= base
	}

	return res
}

// highLowDiv is a bits.Div64 rewrite using uint conventions from this package, i.e.
// each uint64 (hi/lo) stores the number's upper and lower 18 digits.
func highLowDiv(hi, lo, d uint64) (quo, rem uint64) {
	if d == 0 {
		panic("division by zero")
	}
	if d <= hi {
		panic("unsupported denominator")
	}

	// If high part is zero, we can directly return the results.
	if hi == 0 {
		return lo / d, lo % d
	}

	// Left shift d up to maxDigitsPerUint digits.
	s := maxDigitsPerUint - digitsOf(d)
	d, _ = leftShift(d, s)

	const half = halfMaxValuePerUint + 1

	// Splitting d into its lower and upper digits.
	hiD, loD := splitInHalf(d)

	// Getting the 's' left digits of 'lo' and the remaining digits.
	loRaisedDigits, remainingLo := rightShift(lo, maxDigitsPerUint-s)
	// Shifting 'hi' by 's' digits to the left.
	// The overflow is always zero, since digits of 'hi' is less or equal 's', given that 'd' > 'hi'.
	hiShiftedLeft, _ := leftShift(hi, s)
	// upperDigits is the upper digits of the entire number (hi, lo) shifted left by 's' digits.
	upperDigits := hiShiftedLeft + loRaisedDigits

	loShiftedHi, loShiftedLo := splitInHalf(remainingLo)

	hiQuotient := upperDigits / hiD
	hiRemainder := upperDigits - hiQuotient*hiD

	for hiQuotient >= half || hiQuotient*loD > half*hiRemainder+loShiftedHi {
		hiQuotient--
		hiRemainder += hiD
		if hiRemainder >= half {
			break
		}
	}

	un21 := upperDigits*half + loShiftedHi - hiQuotient*d
	q0 := un21 / hiD
	hiRemainder = un21 - q0*hiD

	for q0 >= half || q0*loD > half*hiRemainder+loShiftedLo {
		q0--
		hiRemainder += hiD
		if hiRemainder >= half {
			break
		}
	}

	resultingRem, _ := rightShift(un21*half+loShiftedLo-q0*d, s)

	return hiQuotient*half + q0, resultingRem
}
