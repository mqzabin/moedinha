package moedinha

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var fuzzNumberSeed = generateSeed()

func generateSeed() [][natNumberOfUints]uint64 {

	preset := [][natNumberOfUints]uint64{
		{maxValuePerUint, maxValuePerUint, maxValuePerUint, maxValuePerUint},
		{0, maxValuePerUint, maxValuePerUint, maxValuePerUint},
		{0, 0, maxValuePerUint, maxValuePerUint},
		{0, 0, 0, maxValuePerUint},
		{maxValuePerUint, maxValuePerUint, maxValuePerUint, 0},
		{maxValuePerUint, maxValuePerUint, 0, 0},
		{maxValuePerUint, 0, 0, 0},
		{0, 0, 0, 0},
	}

	forCardinality := maxDigitsPerUint * maxDigitsPerUint * maxDigitsPerUint * maxDigitsPerUint / 16

	seed := make([][natNumberOfUints]uint64, 0, forCardinality+len(preset))
	seed = append(seed, preset...)

	for i := 0; i < maxDigitsPerUint/2; i += 2 {
		for j := 1; j < maxDigitsPerUint/2; j += 2 {
			for k := 0; k < maxDigitsPerUint/2; k += 2 {
				for l := 1; l < maxDigitsPerUint/2; l += 2 {
					a := basePow(i) * uint64((i%2+j%3+k%4+l%5)%10)
					b := basePow(i) * uint64((j%2+k%3+l%4+i%5)%10)
					c := basePow(i) * uint64((k%2+l%3+i%4+j%5)%10)
					d := basePow(i) * uint64((l%2+i%3+j%4+k%5)%10)
					seed = append(seed, [natNumberOfUints]uint64{a, b, c, d})
				}
			}
		}
	}

	return seed
}

func fuzzyUnaryOperationSeedGenerator(f *testing.F) {
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
}

func fuzzyBinaryOperationSeedGenerator(f *testing.F) {
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
}

func fuzzyUnaryOperation(f *testing.F, maxDigits int, fn func(*testing.T, Currency)) {
	f.Helper()

	fuzzyUnaryOperationSeedGenerator(f)

	f.Fuzz(func(t *testing.T, n1, n2, n3, n4 uint64, neg bool) {
		t.Helper()

		curr := seedParser(t, maxDigits, n1, n2, n3, n4, neg)

		fn(t, curr)
	})
}

func fuzzyBinaryOperation(f *testing.F, maxDigits int, fn func(*testing.T, Currency, Currency)) {
	f.Helper()

	fuzzyBinaryOperationSeedGenerator(f)

	f.Fuzz(func(t *testing.T, aN1, aN2, aN3, aN4 uint64, aNeg bool, bN1, bN2, bN3, bN4 uint64, bNeg bool) {
		t.Parallel()

		a := seedParser(t, maxDigits, aN1, aN2, aN3, aN4, aNeg)
		bStr := seedParser(t, maxDigits, bN1, bN2, bN3, bN4, bNeg)

		fn(t, a, bStr)
	})
}

func seedParser(t *testing.T, maxDigits int, n1, n2, n3, n4 uint64, neg bool) Currency {
	t.Helper()

	n1 = n1 % basePow(maxDigitsPerUint)
	if maxDigits < maxDigitsPerUint {
		n1 = n1 % basePow(maxDigits)
	}

	n2 = n2 % maxValuePerUint
	if maxDigits < 2*maxDigitsPerUint {
		n2 = n2 % basePow(maxDigits-maxDigitsPerUint)
	}

	n3 = n3 % maxValuePerUint
	if maxDigits < 3*maxDigitsPerUint {
		n3 = n3 % basePow(maxDigits-2*maxDigitsPerUint)
	}

	n4 = n4 % maxValuePerUint
	if maxDigits < 4*maxDigitsPerUint {
		n4 = n4 % basePow(maxDigits-3*maxDigitsPerUint)
	}

	s := fmt.Sprintf("%018d%018d%018d%018d", n4, n3, n2, n1)
	s = s[:currencyMaxIntegerDigits] + string(currencyDecimalSeparatorSymbol) + s[currencyMaxIntegerDigits:]
	if neg {
		s = "-" + s
	}

	curr, err := NewFromString(s)
	require.NoError(t, err)

	return curr
}

func basePow(n int) uint64 {
	r := uint64(1)
	for i := 0; i < n; i++ {
		r *= base
	}

	return r
}
