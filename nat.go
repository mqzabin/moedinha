package moedinha

import (
	"fmt"
)

const (
	decimalBase = 10
	// natMaxValuePerInt is the greater 999-ish number under 63 bits.
	natMaxValuePerInt = 999999999999999999
	// natMaxDigitsPerInt is the amount of natMaxValuePerInt digits.
	natMaxDigitsPerInt = 18
	// natNumberOfInts stores the amount of uint64 used to represent the currency.
	natNumberOfInts = 3
	// natDigits is the total number of digits that a natural number can have.
	natDigits = natNumberOfInts * natMaxDigitsPerInt
)

type nat struct {
	n1, n2, n3 uint64
}

// newNatFromString v should have natMaxDigitsPerInt*natNumberOfInts length.
func newNatFromString(v [natDigits]byte) (nat, error) {
	// Parsing n3
	var n3Str [natMaxDigitsPerInt]byte
	copy(n3Str[:], v[:natMaxDigitsPerInt])

	n3, err := atoi(n3Str)
	if err != nil {
		return nat{}, fmt.Errorf("error decoding natural number: %w", err)
	}

	// Parsing n2
	var n2Str [natMaxDigitsPerInt]byte
	copy(n2Str[:], v[natMaxDigitsPerInt:2*natMaxDigitsPerInt])

	n2, err := atoi(n2Str)
	if err != nil {
		return nat{}, fmt.Errorf("error decoding natural number: %w", err)
	}

	// Parsing n1
	var n1Str [natMaxDigitsPerInt]byte
	copy(n1Str[:], v[2*natMaxDigitsPerInt:])

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
	copy(str[:natMaxDigitsPerInt], n3[:])

	// Filling n2
	n2 := itoa(n.n2)
	copy(str[natMaxDigitsPerInt:2*natMaxDigitsPerInt], n2[:])

	// Filling n1
	n1 := itoa(n.n1)
	copy(str[2*natMaxDigitsPerInt:], n1[:])

	return str
}

func (n nat) add(v nat) nat {
	r3 := n.n3 + v.n3
	r2 := n.n2 + v.n2
	r1 := n.n1 + v.n1

	r1, r2 = rebalance(r1, r2)
	r2, r3 = rebalance(r2, r3)

	if r3 > natMaxValuePerInt {
		panic("natural number overflow")
	}

	return nat{n1: r1, n2: r2, n3: r3}
}

// difference "n" should always lesser than "v"
func (n nat) difference(v nat) nat {
	vCompl := v.complementOf9()
	sum := n.add(vCompl)

	return nat{
		n1: natMaxValuePerInt - sum.n1,
		n2: natMaxValuePerInt - sum.n2,
		n3: natMaxValuePerInt - sum.n3,
	}
}

func (n nat) complementOf9() nat {
	r3 := natMaxValuePerInt - n.n3
	r2 := natMaxValuePerInt - n.n2
	r1 := natMaxValuePerInt - n.n1

	return nat{n1: r1, n2: r2, n3: r3}
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

func rebalance(src, dest uint64) (newSrc, newDest uint64) {
	if src <= natMaxValuePerInt {
		return src, dest
	}

	dest += src / (natMaxValuePerInt + 1)
	src %= natMaxValuePerInt + 1

	return src, dest
}

// atoi is a fork from strconv.Atoi returning uint64
func atoi(s [natMaxDigitsPerInt]byte) (uint64, error) {
	var n uint64
	for _, ch := range s {
		ch -= '0'
		if ch > 9 {
			return 0, fmt.Errorf("invalid syntax converting string to uint64: rune %c", ch)
		}
		n = n*10 + uint64(ch)
	}

	return n, nil
}

func itoa(v uint64) [natMaxDigitsPerInt]byte {
	var res [natMaxDigitsPerInt]byte
	div := v

	for i := range res {
		digit := byte(div % decimalBase)
		res[natMaxDigitsPerInt-i-1] = zeroRune + digit

		div /= decimalBase
	}

	return res
}
