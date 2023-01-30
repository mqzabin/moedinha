package currency

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func fuzzySeedGenerator(f *testing.F) {
	f.Helper()

	fuzzSeed := [][3]uint64{
		{123456789012345678, 901234567890123456, 789012345678901234},
		{0, 0, 4},
		{0, 4, 0},
		{4, 0, 0},
		{0, 0, 123456789012345678},
		{0, 123456789012345678, 0},
		{123456789012345678, 0, 0},
		{0, 0, 0},
	}

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

func seedParser(t *testing.T, n1, n2, n3 uint64, neg bool) string {
	t.Helper()

	maxValueToNotOverflowSum := uint64(natMaxValuePerInt / 10)

	n1 = n1 % maxValueToNotOverflowSum
	n2 = n2 % maxValueToNotOverflowSum
	n3 = n3 % maxValueToNotOverflowSum

	s := fmt.Sprintf("%018d%018d%018d", n3, n2, n1)
	s = s[:currMaxIntegerDigits] + string(currDecimalPointerSymbol) + s[currMaxIntegerDigits:]
	if neg {
		s = "-" + s
	}

	return s
}

func FuzzAdd(f *testing.F) {
	fuzzySeedGenerator(f)

	f.Fuzz(func(t *testing.T, aN1, aN2, aN3 uint64, aNeg bool, bN1, bN2, bN3 uint64, bNeg bool) {
		aStr := seedParser(t, aN1, aN2, aN3, aNeg)
		t.Log(aStr, " +")

		bStr := seedParser(t, bN1, bN2, bN3, bNeg)
		t.Log(bStr, " =")

		a, err := FromDecimalString(aStr)
		require.NoError(t, err)

		b, err := FromDecimalString(bStr)
		require.NoError(t, err)

		sa, err := decimal.NewFromString(aStr)
		require.NoError(t, err)

		sb, err := decimal.NewFromString(bStr)
		require.NoError(t, err)

		result := a.Add(b)
		t.Log(result.String())
		require.True(t, result.Equal(b.Add(a)))

		sResult := sa.Add(sb)
		t.Log(sResult.String(), " (shopspring)")
		require.Equal(t, sResult.String(), result.String())

		t.Log("pass!")
	})
}

func FuzzSub(f *testing.F) {
	fuzzySeedGenerator(f)

	f.Fuzz(func(t *testing.T, aN1, aN2, aN3 uint64, aNeg bool, bN1, bN2, bN3 uint64, bNeg bool) {
		aStr := seedParser(t, aN1, aN2, aN3, aNeg)
		t.Log(aStr, " -")

		bStr := seedParser(t, bN1, bN2, bN3, bNeg)
		t.Log(bStr, " =")

		a, err := FromDecimalString(aStr)
		require.NoError(t, err)

		b, err := FromDecimalString(bStr)
		require.NoError(t, err)

		sa, err := decimal.NewFromString(aStr)
		require.NoError(t, err)

		sb, err := decimal.NewFromString(bStr)
		require.NoError(t, err)

		result := a.Sub(b)
		t.Log(result.String())

		sResult := sa.Sub(sb)
		t.Log(sResult.String(), " (shopspring)")
		require.Equal(t, sResult.String(), result.String())

		t.Log("pass!")
	})
}

func FuzzGreaterThan(f *testing.F) {
	fuzzySeedGenerator(f)

	f.Fuzz(func(t *testing.T, aN1, aN2, aN3 uint64, aNeg bool, bN1, bN2, bN3 uint64, bNeg bool) {
		aStr := seedParser(t, aN1, aN2, aN3, aNeg)
		t.Log(aStr, " >")

		bStr := seedParser(t, bN1, bN2, bN3, bNeg)
		t.Log(bStr, " =")

		a, err := FromDecimalString(aStr)
		require.NoError(t, err)

		b, err := FromDecimalString(bStr)
		require.NoError(t, err)

		sa, err := decimal.NewFromString(aStr)
		require.NoError(t, err)

		sb, err := decimal.NewFromString(bStr)
		require.NoError(t, err)

		result := a.GreaterThan(b)
		t.Log("Result: ", result)

		sResult := sa.GreaterThan(sb)
		t.Log("Result: ", sResult, " (shospring)")
		require.Equal(t, sResult, result)

		t.Log("pass!")
	})
}

func FuzzLessThan(f *testing.F) {
	fuzzySeedGenerator(f)

	f.Fuzz(func(t *testing.T, aN1, aN2, aN3 uint64, aNeg bool, bN1, bN2, bN3 uint64, bNeg bool) {
		aStr := seedParser(t, aN1, aN2, aN3, aNeg)
		t.Log(aStr, " <")

		bStr := seedParser(t, bN1, bN2, bN3, bNeg)
		t.Log(bStr, " =")

		a, err := FromDecimalString(aStr)
		require.NoError(t, err)

		b, err := FromDecimalString(bStr)
		require.NoError(t, err)

		sa, err := decimal.NewFromString(aStr)
		require.NoError(t, err)

		sb, err := decimal.NewFromString(bStr)
		require.NoError(t, err)

		result := a.LessThan(b)
		t.Log("Result: ", result)

		sResult := sa.LessThan(sb)
		t.Log("Result: ", sResult, " (shospring)")
		require.Equal(t, sResult, result)

		t.Log("pass!")
	})
}

func BenchmarkCurrency(b *testing.B) {

	as := strings.Repeat("7", natMaxDigitsPerInt-1)
	as += strings.Repeat("9", natMaxDigitsPerInt)
	as += strings.Repeat("9", natMaxDigitsPerInt-2)
	as += "80"
	as = as[:len(as)-currDecimalDigits] + "." + as[len(as)-currDecimalDigits:]

	bs := strings.Repeat("0", natMaxDigitsPerInt*natNumberOfInts-2)
	bs += "20"
	bs = bs[:len(bs)-currDecimalDigits] + "." + bs[len(bs)-currDecimalDigits:]

	x, err := FromDecimalString(as)
	require.NoError(b, err)

	y, err := FromDecimalString(bs)
	require.NoError(b, err)

	b.Log(x.String())
	b.Log(y.String())
	b.Log(x.Add(y).String())

	var (
		mCurrency Currency
		sCurrency decimal.Decimal
	)

	b.Run("Add - PoC", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			x, _ := FromDecimalString(as)

			y, _ := FromDecimalString(bs)

			mCurrency = x.Add(y)
		}
	})

	b.Run("Add - shopspring", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			x, _ := decimal.NewFromString(as)

			y, _ := decimal.NewFromString(bs)

			sCurrency = x.Add(y)
		}
	})

	b.Log(mCurrency.String())
	b.Log(sCurrency.String())
}
