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

func (n natural) digits() int {
	for i := 0; i < numberOfUints; i++ {
		if d := digitsOf(n[i]); d > 0 {
			return naturalMaxLen - maxDigitsPerUint*i - (maxDigitsPerUint - d)
		}
	}

	return 0
}

func (n natural) uintsInUse() int {
	for i := 0; i < numberOfUints; i++ {
		if n[i] > 0 {
			return numberOfUints - i
		}
	}

	return 0
}

// rightShiftDigit moves the digits to the right.
// This operation is equivalent to n*base^(-shift).
// The second return is the operation overflow to the right, called "loss".
func (n natural) rightShiftDigit(shift int) (natural, natural) {
	if n.isZero() {
		return natural{}, natural{}
	}

	if shift < 0 {
		return n.leftShiftDigit(-shift)
	}

	if shift > naturalMaxLen {
		panic("invalid shift value passed to right shift digit function")
	}

	shifted, loss := n.rightShiftUint(shift / maxDigitsPerUint)
	shift = shift % maxDigitsPerUint

	if shift == 0 {
		return shifted, loss
	}

	var prevLoss, currLoss uint64

	for i := 0; i < numberOfUints; i++ {
		shifted[i], currLoss = rightShift(shifted[i], shift)
		shifted[i] += prevLoss
		prevLoss = currLoss
	}

	loss, _ = loss.rightShiftDigit(shift)
	loss[0] += prevLoss

	return shifted, loss
}

// leftShiftUint moves the uint64 components of the natural number to left.
// This operation is equivalent to n*base^shift.
// The second return is the operation overflow.
func (n natural) leftShiftDigit(shift int) (natural, natural) {
	if n.isZero() {
		return natural{}, natural{}
	}

	if shift < 0 {
		return n.rightShiftDigit(-shift)
	}

	if shift > naturalMaxLen {
		panic("invalid shift value passed to left shift digit function")
	}

	shifted, overflow := n.leftShiftUint(shift / maxDigitsPerUint)
	shift = shift % maxDigitsPerUint

	if shift == 0 {
		return shifted, overflow
	}

	var prevOverflow, currOverflow uint64

	for i := numberOfUints - 1; i >= 0; i-- {
		shifted[i], currOverflow = leftShift(shifted[i], shift)
		shifted[i] += prevOverflow
		prevOverflow = currOverflow
	}

	overflow, _ = overflow.leftShiftDigit(shift)
	overflow[numberOfUints-1] += prevOverflow

	return shifted, overflow
}

// rightShiftUint moves the uint64 components of the natural number to right.
// This operation is equivalent to n^(-shift*maxDigitsPerUint).
// The second return is the operation overflow to the right, called "loss".
func (n natural) rightShiftUint(shift int) (natural, natural) {
	if n.isZero() {
		return natural{}, natural{}
	}

	if shift == 0 {
		return n, natural{}
	}

	if shift < 0 {
		return n.leftShiftUint(-shift)
	}

	if shift >= 2*numberOfUints {
		panic("invalid shift value passed to right shift uint function")
	}

	var shifted, loss natural

	for i := 0; i < numberOfUints; i++ {
		if i+shift < numberOfUints {
			shifted[i+shift] = n[i]
			continue
		}
		loss[i+shift-numberOfUints] = n[i]
	}

	return shifted, loss
}

// leftShiftUint moves the components of the natural number to left.
// This operation is equivalent to n^(shift*maxDigitsPerUint).
// The second return is the operation overflow.
func (n natural) leftShiftUint(shift int) (natural, natural) {
	if n.isZero() {
		return natural{}, natural{}
	}

	if shift == 0 {
		return n, natural{}
	}

	if shift < 0 {
		return n.rightShiftUint(-shift)
	}

	if shift >= 2*numberOfUints {
		panic("invalid shift value passed to left shift uint function")
	}

	var shifted, overflow natural

	for i := 0; i < numberOfUints; i++ {
		if i-shift < 0 {
			overflow[numberOfUints+i-shift] = n[i]
			continue
		}

		shifted[i-shift] = n[i]
	}

	return shifted, overflow
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

		padded, paddingOverflow := mr.leftShiftUint(numberOfUints - i - 1)

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

func (n natural) div(v natural) (natural, natural) {
	if v.isZero() {
		panic("division by zero")
	}

	if n.isZero() {
		return natural{}, natural{}
	}

	if n.lessThan(v) {
		return natural{}, natural{}
	}

	if n.equal(v) {
		return naturalOne, natural{}
	}

	vUintsInUse := v.uintsInUse()
	nUintsInUse := n.uintsInUse()

	// Since n > v, we could assume that nUintsInUse >= vUintsInUse.

	// Both have 1 uint64, so it's a simple uint64 division.
	if vUintsInUse == 1 && nUintsInUse == 1 {
		var quotient, remainder natural
		quotient[numberOfUints-1] = n[numberOfUints-1] / v[numberOfUints-1]
		remainder[numberOfUints-1] = n[numberOfUints-1] % v[numberOfUints-1]

		return quotient, remainder
	}

	// Denominator has 1 uint64, so divByUint could be used.
	if vUintsInUse == 1 {
		quotient, remUint := n.divByUint(v[numberOfUints-1])

		var remainder natural
		remainder[numberOfUints-1] = remUint

		return quotient, remainder
	}

	var quotient natural

	dividend := n
	for i := 0; i < nUintsInUse; i++ {
		var digit uint64
		digit, dividend = longDivisionIteration(dividend, v)

		// Inserting new computed digit to the quotient.
		quotient, _ = quotient.leftShiftUint(1)
		quotient[numberOfUints-1] = digit
	}

	// The carried dividend is the final remainder.
	return quotient, dividend
}

// longDivisionIteration searches for the minimum digits of 'n' that compose a number greater than 'd',
// then divides the number found by 'd', returning the quotient, and the remainder of 'n'.
// Due to this operation, the quotient will never be greater than 1 uint64.
// The remainder is computed by `n - quotient*d*B^(shiftedDigits)`.
//
// This function is mainly used by natural.div method in the long division algorithm.
func longDivisionIteration(n, d natural) (uint64, natural) {
	nUints := n.uintsInUse()
	dUints := d.uintsInUse()

	var hiN, loN uint64
	loN = n[numberOfUints-nUints]

	var hiD uint64
	hiD = d[numberOfUints-dUints]

	if hiD > loN {
		hiN = loN
		loN = n[numberOfUints-nUints-1]
	}

	quoDigit, _ := highLowDiv(hiN, loN, hiD)

	quoMul, _ := d.mulByUint64(quoDigit)
	quoMul, _ = quoMul.leftShiftUint(nUints - dUints)

	return quoDigit, n.sub(quoMul)
}

func (n natural) divByUint(v uint64) (natural, uint64) {
	var (
		remainderSum uint64
		result       natural
	)

	// Partially dividing each digit of 'n', since
	// (B^n + ... + B^0)/v = (B^n/v) + ... + (B^0/v)
	for i := 0; i < numberOfUints; i++ {
		var partialN, partialResult natural

		// Use only the i-th digit of 'n'.
		// 0 0 0 ... n[i] ... 0 0 0
		partialN[i] = n[i]

		var carriedRemainder uint64
		// Dividing each digit of partialN.
		for j := i; j < numberOfUints; j++ {
			partialResult[j], carriedRemainder = highLowDiv(carriedRemainder, partialN[j], v)
		}

		result = result.add(partialResult)
		remainderSum += carriedRemainder
	}

	// Creating a natural with the least significant digit equal to 'remainderSum / v'.
	// Note that 'remainderSum / v' will never overflow maxValuePerUint.
	var remainderDivision natural
	remainderDivision[numberOfUints-1] = remainderSum / v

	result = result.add(remainderDivision)

	return result, remainderSum % v
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
