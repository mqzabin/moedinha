package moedinha

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func FuzzAdd(f *testing.F) {
	fuzzyBinaryOperation(f, natDigits-1, func(t *testing.T, a, b Currency) {
		sa, err := decimal.NewFromString(a.String())
		require.NoError(t, err)

		sb, err := decimal.NewFromString(b.String())
		require.NoError(t, err)

		result := a.Add(b)

		require.True(t, result.Equal(b.Add(a)))

		sResult := sa.Add(sb)

		require.Equal(t, sResult.String(), result.String())

	})
}

func FuzzSub(f *testing.F) {
	fuzzyBinaryOperation(f, natDigits-1, func(t *testing.T, a, b Currency) {

		sa, err := decimal.NewFromString(a.String())
		require.NoError(t, err)

		sb, err := decimal.NewFromString(b.String())
		require.NoError(t, err)

		result := a.Sub(b)

		sResult := sa.Sub(sb)

		require.Equal(t, sResult.String(), result.String())

	})
}

func FuzzMul(f *testing.F) {
	fuzzyBinaryOperation(f, natDigits/2, func(t *testing.T, a, b Currency) {

		sa, err := decimal.NewFromString(a.String())
		require.NoError(t, err)

		sb, err := decimal.NewFromString(b.String())
		require.NoError(t, err)

		if as, bs := sa.StringFixed(currencyDecimalDigits), sb.StringFixed(currencyDecimalDigits); len(as)+len(bs)-4 > natDigits {
			t.Skip(as, bs)
		}

		result := a.Mul(b)

		sResult := sa.Mul(sb).Truncate(currencyDecimalDigits)

		require.Equalf(t, sResult.String(), result.String(), "a: %s, b: %s", a.String(), b.String())
	})
}

func FuzzGreaterThan(f *testing.F) {
	fuzzyBinaryOperation(f, natDigits, func(t *testing.T, a, b Currency) {

		sa, err := decimal.NewFromString(a.String())
		require.NoError(t, err)

		sb, err := decimal.NewFromString(b.String())
		require.NoError(t, err)

		result := a.GreaterThan(b)

		sResult := sa.GreaterThan(sb)

		require.Equal(t, sResult, result)
	})
}

func FuzzLessThan(f *testing.F) {
	fuzzyBinaryOperation(f, natDigits, func(t *testing.T, a, b Currency) {

		sa, err := decimal.NewFromString(a.String())
		require.NoError(t, err)

		sb, err := decimal.NewFromString(b.String())
		require.NoError(t, err)

		result := a.LessThan(b)

		sResult := sa.LessThan(sb)

		require.Equal(t, sResult, result)
	})
}

func FuzzEqual(f *testing.F) {
	fuzzyBinaryOperation(f, natDigits, func(t *testing.T, a, b Currency) {

		sa, err := decimal.NewFromString(a.String())
		require.NoError(t, err)

		sb, err := decimal.NewFromString(b.String())
		require.NoError(t, err)

		result := a.Equal(b)

		sResult := sa.Equal(sb)

		require.Equal(t, sResult, result)
	})
}

func FuzzIsZero(f *testing.F) {
	fuzzyUnaryOperation(f, natDigits, func(t *testing.T, a Currency) {
		sa, err := decimal.NewFromString(a.String())
		require.NoError(t, err)

		result := a.IsZero()

		sResult := sa.IsZero()

		require.Equal(t, sResult, result)
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
	aStr := "10000000000000000010000000000000000010000.00000000000001"
	bStr := "0.999999999999999999"

	var (
		mCurrency Currency
		sCurrency decimal.Decimal
	)

	b.Run("moedinha", func(b *testing.B) {
		x, _ := NewFromString(aStr)

		y, _ := NewFromString(bStr)

		for i := 0; i < b.N; i++ {
			mCurrency = x.Mul(y)
		}
	})

	b.Run("shopspring", func(b *testing.B) {
		x, _ := decimal.NewFromString(aStr)

		y, _ := decimal.NewFromString(bStr)

		for i := 0; i < b.N; i++ {
			sCurrency = x.Mul(y)
		}
	})

	b.Log(mCurrency.String())
	b.Log(sCurrency.String())
}

func BenchmarkMul(b *testing.B) {
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
