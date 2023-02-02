package moedinha

import "fmt"

const (
	// base is the base value used by the library.
	base = 10
	// maxValuePerUint is the greater 999-ish number under 63 bits.
	maxValuePerUint = 999999999999999999
	// maxDigitsPerUint is the amount of maxValuePerUint digits.
	maxDigitsPerUint = 18
	// halfMaxValuePerUint is half of maxValuePerUint digits.
	halfMaxValuePerUint = 999999999
)

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
