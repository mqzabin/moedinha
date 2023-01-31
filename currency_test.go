package moedinha

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func FuzzAdd(f *testing.F) {
	fuzzyBinaryOperation(f, func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " +")
		t.Log(bStr, " =")

		a, err := NewFromString(aStr)
		require.NoError(t, err)

		b, err := NewFromString(bStr)
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
	fuzzyBinaryOperation(f, func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " -")
		t.Log(bStr, " =")

		a, err := NewFromString(aStr)
		require.NoError(t, err)

		b, err := NewFromString(bStr)
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
	fuzzyBinaryOperation(f, func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " >")
		t.Log(bStr, " =")

		a, err := NewFromString(aStr)
		require.NoError(t, err)

		b, err := NewFromString(bStr)
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
	fuzzyBinaryOperation(f, func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " <")
		t.Log(bStr, " =")

		a, err := NewFromString(aStr)
		require.NoError(t, err)

		b, err := NewFromString(bStr)
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

func FuzzEqual(f *testing.F) {
	fuzzyBinaryOperation(f, func(t *testing.T, aStr, bStr string) {
		t.Log(aStr, " ==")
		t.Log(bStr, " =")

		a, err := NewFromString(aStr)
		require.NoError(t, err)

		b, err := NewFromString(bStr)
		require.NoError(t, err)

		sa, err := decimal.NewFromString(aStr)
		require.NoError(t, err)

		sb, err := decimal.NewFromString(bStr)
		require.NoError(t, err)

		result := a.Equal(b)
		t.Log("Result: ", result)

		sResult := sa.Equal(sb)
		t.Log("Result: ", sResult, " (shospring)")
		require.Equal(t, sResult, result)

		t.Log("pass!")
	})
}

func FuzzIsZero(f *testing.F) {
	fuzzyUnaryOperation(f, func(t *testing.T, str string) {
		t.Log(str, " is zero?")

		a, err := NewFromString(str)
		require.NoError(t, err)

		sa, err := decimal.NewFromString(str)
		require.NoError(t, err)

		result := a.IsZero()
		t.Log("Result: ", result)

		sResult := sa.IsZero()
		t.Log("Result: ", sResult, " (shospring)")
		require.Equal(t, sResult, result)

		t.Log("pass!")
	})
}

func BenchmarkNewFromString(b *testing.B) {
	aStr := "8901234567890124190123456789012345612345678.9012345678"

	var (
		mCurrency Currency
		sCurrency decimal.Decimal
	)

	b.Run("moedinha", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			x, _ := NewFromString(aStr)

			mCurrency = x
		}
	})

	b.Run("shopspring", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x, _ := decimal.NewFromString(aStr)

			sCurrency = x
		}
	})

	b.Log(mCurrency.String())
	b.Log(sCurrency.String())
}

func BenchmarkString(b *testing.B) {
	aStr := "8901234567890124190123456789012345612345678.9012345678"

	var (
		mCurrency string
		sCurrency string
	)

	b.Run("moedinha", func(b *testing.B) {
		x, _ := NewFromString(aStr)

		for i := 0; i < b.N; i++ {
			mCurrency = x.String()
		}
	})

	b.Run("shopspring", func(b *testing.B) {
		x, _ := decimal.NewFromString(aStr)

		for i := 0; i < b.N; i++ {

			sCurrency = x.String()
		}
	})

	b.Log(mCurrency)
	b.Log(sCurrency)
}

func BenchmarkAdd(b *testing.B) {
	aStr := "8901234567890124190123456789012345612345678.9012345678"
	bStr := "2345678901234567500000000000000000000000000"

	var (
		mCurrency Currency
		sCurrency decimal.Decimal
	)

	b.Run("moedinha", func(b *testing.B) {
		x, _ := NewFromString(aStr)

		y, _ := NewFromString(bStr)

		for i := 0; i < b.N; i++ {
			mCurrency = x.Add(y)
		}
	})

	b.Run("shopspring", func(b *testing.B) {
		x, _ := decimal.NewFromString(aStr)

		y, _ := decimal.NewFromString(bStr)

		for i := 0; i < b.N; i++ {
			sCurrency = x.Add(y)
		}
	})

	b.Log(mCurrency.String())
	b.Log(sCurrency.String())
}

func BenchmarkSub(b *testing.B) {
	aStr := "8901234567890124190123456789012345612345678.9012345678"
	bStr := "2345678901234567500000000000000000000000000"

	var (
		mCurrency Currency
		sCurrency decimal.Decimal
	)

	b.Run("moedinha", func(b *testing.B) {
		x, _ := NewFromString(aStr)

		y, _ := NewFromString(bStr)

		for i := 0; i < b.N; i++ {
			mCurrency = x.Sub(y)
		}
	})

	b.Run("shopspring", func(b *testing.B) {
		x, _ := decimal.NewFromString(aStr)

		y, _ := decimal.NewFromString(bStr)

		for i := 0; i < b.N; i++ {
			sCurrency = x.Sub(y)
		}
	})

	b.Log(mCurrency.String())
	b.Log(sCurrency.String())
}
