package moedinha

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var fuzzSeed = [][3]uint64{
	{123456789012345678, 901234567890123456, 789012345678901234},
	{0, 0, 4},
	{0, 4, 0},
	{4, 0, 0},
	{0, 0, 123456789012345678},
	{0, 123456789012345678, 0},
	{123456789012345678, 0, 0},
	{0, 0, 0},
}

func fuzzyUnaryOperationSeedGenerator(f *testing.F) {
	f.Helper()

	for a := 0; a < 2; a++ {
		aNeg := a == 0
		for i := 0; i < len(fuzzSeed); i++ {
			f.Add(
				fuzzSeed[i][0], fuzzSeed[i][1], fuzzSeed[i][2],
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
			for i := 0; i < len(fuzzSeed); i++ {
				for j := 0; j < len(fuzzSeed); j++ {
					f.Add(
						fuzzSeed[i][0], fuzzSeed[i][1], fuzzSeed[i][2], // first number
						aNeg,
						fuzzSeed[j][0], fuzzSeed[j][1], fuzzSeed[j][2], // second number
						bNeg,
					)
				}
			}
		}
	}
}

func fuzzyUnaryOperation(f *testing.F, fn func(*testing.T, string)) {
	f.Helper()

	fuzzyUnaryOperationSeedGenerator(f)

	f.Fuzz(func(t *testing.T, n1, n2, n3 uint64, neg bool) {
		str := seedParser(t, n1, n2, n3, neg)

		fn(t, str)
	})
}

func fuzzyBinaryOperation(f *testing.F, fn func(*testing.T, string, string)) {
	f.Helper()

	fuzzyBinaryOperationSeedGenerator(f)

	f.Fuzz(func(t *testing.T, aN1, aN2, aN3 uint64, aNeg bool, bN1, bN2, bN3 uint64, bNeg bool) {
		aStr := seedParser(t, aN1, aN2, aN3, aNeg)
		bStr := seedParser(t, bN1, bN2, bN3, bNeg)

		a, err := NewFromString(aStr)
		require.NoError(t, err)

		b, err := NewFromString(bStr)
		require.NoError(t, err)

		fn(t, a.String(), b.String())
		//fn(t, aStr, bStr)
	})
}

func seedParser(t *testing.T, n1, n2, n3 uint64, neg bool) string {
	t.Helper()

	maxValueToNotOverflowSum := uint64(natMaxValuePerInt / 10)

	n1 = n1 % natMaxValuePerInt
	n2 = n2 % natMaxValuePerInt
	n3 = n3 % maxValueToNotOverflowSum

	s := fmt.Sprintf("%018d%018d%018d", n3, n2, n1)
	s = s[:currencyMaxIntegerDigits] + string(currencyDecimalSeparatorSymbol) + s[currencyMaxIntegerDigits:]
	if neg {
		s = "-" + s
	}

	return s
}
