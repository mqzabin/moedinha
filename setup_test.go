package moedinha

import (
	"strings"
	"testing"
)

var fuzzNumberSeed = generateSeed()

type fuzzSeed struct {
	n   [numberOfUints]uint64
	neg bool
}

func (fs fuzzSeed) string(truncateTo int) string {
	var result string

	for i := range fs.n {
		fs.n[i], _ = rebalance(fs.n[i], 0)
		nStr := itoa(fs.n[i])
		result += string(nStr[:])
	}

	result = result[:truncateTo]
	if truncateTo < currencyDecimalDigits {
		result = strings.Repeat(string(zeroRune), currencyDecimalDigits-truncateTo) + result
		truncateTo = currencyDecimalDigits
	}

	result = result[:truncateTo-currencyDecimalDigits] + string(currencyDecimalSeparatorSymbol) + result[truncateTo-currencyDecimalDigits:]

	result = strings.TrimLeft(result, string(zeroRune))

	if result[0] == currencyDecimalSeparatorSymbol {
		result = string(zeroRune) + result
	}

	if fs.neg {
		result = "-" + result
	}

	return result
}

func generateMaxArray() [numberOfUints]uint64 {
	var seed [numberOfUints]uint64
	for i := range seed {
		seed[i] = maxValuePerUint
	}

	return seed
}

func generateSeed() [][numberOfUints]uint64 {
	seeds := make([][numberOfUints]uint64, 0, 2*numberOfUints)

	// Upper triangular matrix
	for i := 0; i < numberOfUints; i++ {
		var seed [numberOfUints]uint64

		for j := i; j < numberOfUints; j++ {
			seed[j] = maxValuePerUint
		}

		seeds = append(seeds, seed)
	}

	// Lower triangular matrix
	for i := 0; i < numberOfUints; i++ {
		seed := generateMaxArray()

		for j := i; j < numberOfUints; j++ {
			seed[j] = 0
		}

		seeds = append(seeds, seed)
	}

	return seeds
}

func fuzzyUnaryOperation(f *testing.F, fn func(*testing.T, fuzzSeed)) {
	f.Helper()

	for a := 0; a < 2; a++ {
		aNeg := a == 0
		for i := 0; i < len(fuzzNumberSeed); i++ {
			f.Add(
				fuzzNumberSeed[i][0], fuzzNumberSeed[i][1], fuzzNumberSeed[i][2], fuzzNumberSeed[i][3],
				aNeg,
			)
		}
	}

	f.Fuzz(func(t *testing.T, n1, n2, n3, n4 uint64, neg bool) {
		t.Helper()

		a := fuzzSeed{
			n:   [numberOfUints]uint64{n1, n2, n3, n4},
			neg: neg,
		}

		fn(t, a)
	})
}

func fuzzyBinaryOperation(f *testing.F, fn func(*testing.T, fuzzSeed, fuzzSeed)) {
	f.Helper()

	for a := 0; a < 2; a++ {
		aNeg := a == 0
		for b := 0; b < 2; b++ {
			bNeg := b == 0
			for i := 0; i < len(fuzzNumberSeed); i++ {
				for j := 0; j < len(fuzzNumberSeed); j++ {
					f.Add(
						fuzzNumberSeed[i][0], fuzzNumberSeed[i][1], fuzzNumberSeed[i][2], fuzzNumberSeed[i][3], // first number
						aNeg,
						fuzzNumberSeed[j][0], fuzzNumberSeed[j][1], fuzzNumberSeed[j][2], fuzzNumberSeed[j][3], // second number
						bNeg,
					)
				}
			}
		}
	}

	f.Fuzz(func(t *testing.T, aN1, aN2, aN3, aN4 uint64, aNeg bool, bN1, bN2, bN3, bN4 uint64, bNeg bool) {
		t.Parallel()

		a := fuzzSeed{
			n:   [numberOfUints]uint64{aN1, aN2, aN3, aN4},
			neg: aNeg,
		}

		b := fuzzSeed{
			n:   [numberOfUints]uint64{bN1, bN2, bN3, bN4},
			neg: bNeg,
		}

		fn(t, a, b)
	})
}
