package moedinha

import (
	"fmt"
)

const (
	// natNumberOfInts stores the amount of uint64 used to represent the currency.
	natNumberOfInts = 3
	// natDigits is the total number of digits that a natural number can have.
	natDigits = natNumberOfInts * maxDigitsPerUint
)

type nat struct {
	n1, n2, n3 uint64
}

// newNatFromString v should have maxDigitsPerUint*natNumberOfInts length.
func newNatFromString(v [natDigits]byte) (nat, error) {
	// Parsing n3
	var n3Str [maxDigitsPerUint]byte
	copy(n3Str[:], v[:maxDigitsPerUint])

	n3, err := atoi(n3Str)
	if err != nil {
		return nat{}, fmt.Errorf("error decoding natural number: %w", err)
	}

	// Parsing n2
	var n2Str [maxDigitsPerUint]byte
	copy(n2Str[:], v[maxDigitsPerUint:2*maxDigitsPerUint])

	n2, err := atoi(n2Str)
	if err != nil {
		return nat{}, fmt.Errorf("error decoding natural number: %w", err)
	}

	// Parsing n1
	var n1Str [maxDigitsPerUint]byte
	copy(n1Str[:], v[2*maxDigitsPerUint:])

	n1, err := atoi(n1Str)
	if err != nil {
		return nat{}, fmt.Errorf("error decoding natural number: %w", err)
	}

	return nat{
		n1: n1,
		n2: n2,
		n3: n3,
	}, nil
}

func (n nat) string() [natDigits]byte {
	var str [natDigits]byte

	// Filling n3
	n3 := itoa(n.n3)
	copy(str[:maxDigitsPerUint], n3[:])

	// Filling n2
	n2 := itoa(n.n2)
	copy(str[maxDigitsPerUint:2*maxDigitsPerUint], n2[:])

	// Filling n1
	n1 := itoa(n.n1)
	copy(str[2*maxDigitsPerUint:], n1[:])

	return str
}

func (n nat) add(v nat) nat {
	r3 := n.n3 + v.n3
	r2 := n.n2 + v.n2
	r1 := n.n1 + v.n1

	r1, r2 = rebalance(r1, r2)
	r2, r3 = rebalance(r2, r3)

	if r3 > maxValuePerUint {
		panic("natural number overflow")
	}

	return nat{n1: r1, n2: r2, n3: r3}
}

// difference "n" should always lesser than "v"
func (n nat) difference(v nat) nat {
	vCompl := v.complementOf9()
	sum := n.add(vCompl)

	return nat{
		n1: maxValuePerUint - sum.n1,
		n2: maxValuePerUint - sum.n2,
		n3: maxValuePerUint - sum.n3,
	}
}

func (n nat) complementOf9() nat {
	r3 := maxValuePerUint - n.n3
	r2 := maxValuePerUint - n.n2
	r1 := maxValuePerUint - n.n1

	return nat{n1: r1, n2: r2, n3: r3}
}

func (n nat) multiply(v nat) (nat, nat) {
	r1, o1 := n.multiplyByUint(v.n1)
	r2, o2 := n.multiplyByUint(v.n2)
	r3, o3 := n.multiplyByUint(v.n3)

	r := r1.add(r2).add(r3)

	return r, nat{
		n1: o1,
		n2: o2,
		n3: o3,
	}
}

func (n nat) multiplyByUint(x uint64) (nat, uint64) {
	r1, o1 := multiplyUint(n.n1, x)
	r2, o2 := multiplyUint(n.n2, x)
	r3, o3 := multiplyUint(n.n3, x)

	r2 += o1
	r3 += o2

	r1, r2 = rebalance(r1, r2)
	r2, r3 = rebalance(r2, r3)
	r3, o3 = rebalance(r3, o3)

	return nat{
		n1: r1,
		n2: r2,
		n3: r3,
	}, o3
}

func (n nat) isZero() bool {
	return n.n1 == 0 && n.n2 == 0 && n.n3 == 0
}

func (n nat) equal(v nat) bool {
	return n.n1 == v.n1 && n.n2 == v.n2 && n.n3 == v.n3
}

func (n nat) greaterThan(v nat) bool {
	// TODO: Simplify this to a single boolean statement.

	if n.n3 > v.n3 {
		return true
	}

	if n.n3 < v.n3 {
		return false
	}

	// n3 are equal

	if n.n2 > v.n2 {
		return true
	}

	if n.n2 < v.n2 {
		return false
	}

	// n2 are equal

	if n.n1 > v.n1 {
		return true
	}

	if n.n1 < v.n1 {
		return false
	}

	// are equal
	return false
}

func (n nat) lessThan(v nat) bool {
	// TODO: Simplify this to a single boolean statement.

	if n.n3 < v.n3 {
		return true
	}

	if n.n3 > v.n3 {
		return false
	}

	// n3 are equal

	if n.n2 < v.n2 {
		return true
	}

	if n.n2 > v.n2 {
		return false
	}

	// n2 are equal

	if n.n1 < v.n1 {
		return true
	}

	if n.n1 > v.n1 {
		return false
	}

	// are equal
	return false
}
