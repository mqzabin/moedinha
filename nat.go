package moedinha

import (
	"fmt"
)

// naturalMaxLen is the max length of a natural number string .
const naturalMaxLen = numberOfUints * maxDigitsPerUint

// natural represents a natural number.
type natural [numberOfUints]uint64

func newNatFromString(v [naturalMaxLen]byte) (natural, error) {
	var n natural

	for i := 0; i < numberOfUints; i++ {
		var str [maxDigitsPerUint]byte
		copy(str[:], v[i*maxDigitsPerUint:(i+1)*maxDigitsPerUint])

		c, err := atoi(str)
		if err != nil {
			return natural{}, fmt.Errorf("error decoding natural number: %w", err)
		}

		n[i] = c
	}

	return n, nil
}

func (n natural) string() [naturalMaxLen]byte {
	var str [naturalMaxLen]byte

	for i := 0; i < numberOfUints; i++ {
		c := itoa(n[i])
		copy(str[i*maxDigitsPerUint:(i+1)*maxDigitsPerUint], c[:])
	}

	return str
}

// add sums two natural numbers. This operation panics on overflow.
func (n natural) add(v natural) natural {
	var result natural

	for i := numberOfUints - 1; i >= 0; i-- {
		result[i] = n[i] + v[i]

		if i != numberOfUints-1 {
			result[i+1], result[i] = rebalance(result[i+1], result[i])
		}
	}

	var over uint64
	result[0], over = rebalance(result[0], over)

	if over > 0 {
		panic(fmt.Sprintf("natural number overflow: %s + %s", n.string(), v.string()))
	}

	return result
}

// padRight moves the components of the natural number to right.
// This operation is equivalent to n^(-padding*maxDigitsPerUint).
// The second return is the operation overflow to the right, called "loss".
func (n natural) padRight(padding int) (natural, natural) {
	var padded, loss natural

	if padding == 0 {
		return n, natural{}
	}

	if padding < 0 {
		return n.padLeft(-padding)
	}

	if padding >= 2*numberOfUints {
		panic("invalid padding value passed to padRight")
	}

	for i := 0; i < numberOfUints; i++ {
		if i+padding < numberOfUints {
			padded[i+padding] = n[i]
			continue
		}
		loss[i+padding-numberOfUints] = n[i]
	}

	return padded, loss
}

// padLeft moves the components of the natural number to left.
// This operation is equivalent to n^(padding*maxDigitsPerUint).
// The second return is the operation overflow.
func (n natural) padLeft(padding int) (natural, natural) {
	var padded, overflow natural

	if padding == 0 {
		return n, natural{}
	}

	if padding < 0 {
		return n.padRight(-padding)
	}

	if padding >= 2*numberOfUints {
		panic("invalid padding value passed to padLeft")
	}

	for i := 0; i < numberOfUints; i++ {
		if i-padding < 0 {
			overflow[numberOfUints+i-padding] = n[i]
			continue
		}

		padded[i-padding] = n[i]
	}

	return padded, overflow
}

// sub calculates the subtraction: n - v.
// "v" should be lesser than "n" to not overflow the natural domain.
func (n natural) sub(v natural) natural {
	vCompl := n.complement()
	sum := v.add(vCompl)

	var result natural
	for i := 0; i < numberOfUints; i++ {
		result[i] = maxValuePerUint - sum[i]
	}

	return result
}

// complement calculates the complement of n.
// Complement is basically 999...(# of digits) - n.
func (n natural) complement() natural {
	var result natural

	for i := 0; i < numberOfUints; i++ {
		result[i] = maxValuePerUint - n[i]
	}

	return result
}

// mul multiplies two natural numbers.
// The first return is the result, and the second return is the overflow
// of the operation, if any.
func (n natural) mul(v natural) (natural, natural) {
	var result, overflow natural

	for i := 0; i < numberOfUints; i++ {
		mr, mo := n.mulByUint64(v[i])

		padded, paddingOverflow := mr.padLeft(numberOfUints - i - 1)

		overflow[i] += mo
		overflow = overflow.add(paddingOverflow)
		result = result.add(padded)
	}

	return result, overflow
}

// mulByUint64 multiplies a natural number by an uint64.
// The first return is the result, and the second return is the overflow
// of the operation, if any.
func (n natural) mulByUint64(x uint64) (natural, uint64) {
	var result, overflow natural

	for i := numberOfUints - 1; i >= 0; i-- {
		result[i], overflow[i] = multiplyUint(n[i], x)

		if i != numberOfUints-1 {
			result[i] += overflow[i+1]
			result[i+1], result[i] = rebalance(result[i+1], result[i])
		}
	}

	return result, overflow[0]
}

func (n natural) isZero() bool {
	return n == natural{}
}

func (n natural) equal(v natural) bool {
	return n == v
}

func (n natural) greaterThan(v natural) bool {
	return n.greaterOnEqual(v, false)
}

func (n natural) greaterThanOrEqual(v natural) bool {
	return n.greaterOnEqual(v, true)
}

func (n natural) greaterOnEqual(v natural, onEqual bool) bool {
	for i := 0; i < numberOfUints; i++ {
		if n[i] > v[i] {
			return true
		}

		if n[i] < v[i] {
			return false
		}
	}

	// are equal
	return onEqual
}

func (n natural) lessThan(v natural) bool {
	return n.lessOnEqual(v, false)
}

func (n natural) lessThanOrEqual(v natural) bool {
	return n.lessOnEqual(v, false)
}

func (n natural) lessOnEqual(v natural, onEqual bool) bool {
	for i := 0; i < numberOfUints; i++ {
		if n[i] < v[i] {
			return true
		}

		if n[i] > v[i] {
			return false
		}
	}

	// are equal
	return onEqual
}
