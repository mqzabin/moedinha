package currency

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func FuzzAdd(f *testing.F) {
	fuzzyBinaryOpSeedGenerator(f)

	f.Fuzz(fuzzyBinaryOpWrapper(func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " +")
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
	}))
}

func FuzzSub(f *testing.F) {
	fuzzyBinaryOpSeedGenerator(f)

	f.Fuzz(fuzzyBinaryOpWrapper(func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " -")
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
	}))
}

func FuzzGreaterThan(f *testing.F) {
	fuzzyBinaryOpSeedGenerator(f)

	f.Fuzz(fuzzyBinaryOpWrapper(func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " >")
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
	}))
}

func FuzzLessThan(f *testing.F) {
	fuzzyBinaryOpSeedGenerator(f)

	f.Fuzz(fuzzyBinaryOpWrapper(func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " <")
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
	}))
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
