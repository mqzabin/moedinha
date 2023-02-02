package moedinha

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func FuzzBinaryOperations(f *testing.F) {
	fuzzyBinaryOperation(f, func(t *testing.T, seedA, seedB fuzzSeed) {
		t.Run("Add + Sub", func(t *testing.T) {
			// sum could lead to a +1 increase in number of digits.
			truncateToAvoidOverflow := naturalMaxLen - 1

			aStr := seedA.String(truncateToAvoidOverflow)
			bStr := seedB.String(truncateToAvoidOverflow)

			a, err := NewFromString(aStr)
			require.NoError(t, err)

			b, err := NewFromString(bStr)
			require.NoError(t, err)

			sa, err := decimal.NewFromString(aStr)
			require.NoError(t, err)

			sb, err := decimal.NewFromString(bStr)
			require.NoError(t, err)

			addResult := a.Add(b)
			sAddResult := sa.Add(sb)

			require.Equal(t, sAddResult.String(), addResult.String())

			subResult := a.Sub(b)
			sSubResult := sa.Sub(sb)

			require.Equal(t, sSubResult.String(), subResult.String())
		})

		t.Run("Mul", func(t *testing.T) {
			// multiplication could double the number of digits.
			truncateToAvoidOverflow := naturalMaxLen / 2

			aStr := seedA.String(truncateToAvoidOverflow)
			bStr := seedB.String(truncateToAvoidOverflow)

			a, err := NewFromString(aStr)
			require.NoError(t, err)

			b, err := NewFromString(bStr)
			require.NoError(t, err)

			sa, err := decimal.NewFromString(aStr)
			require.NoError(t, err)

			sb, err := decimal.NewFromString(bStr)
			require.NoError(t, err)

			mulResult := a.Mul(b)
			sMulResult := sa.Mul(sb).Truncate(currencyDecimalDigits)

			require.Equal(t, sMulResult.String(), mulResult.String())
		})

		t.Run("Comparisons", func(t *testing.T) {
			// no overflow can occur.
			truncateToAvoidOverflow := naturalMaxLen

			aStr := seedA.String(truncateToAvoidOverflow)
			bStr := seedB.String(truncateToAvoidOverflow)

			a, err := NewFromString(aStr)
			require.NoError(t, err)

			b, err := NewFromString(bStr)
			require.NoError(t, err)

			sa, err := decimal.NewFromString(aStr)
			require.NoError(t, err)

			sb, err := decimal.NewFromString(bStr)
			require.NoError(t, err)

			require.Equal(t, sa.Equal(sb), a.Equal(b))
			require.Equal(t, sa.GreaterThan(sb), a.GreaterThan(b))
			require.Equal(t, sa.GreaterThanOrEqual(sb), a.GreaterOrEqualThan(b))
			require.Equal(t, sa.LessThan(sb), a.LessThan(b))
			require.Equal(t, sa.LessThanOrEqual(sb), a.LessOrEqualThan(b))
		})
	})
}

func FuzzUnaryOperations(f *testing.F) {
	fuzzyUnaryOperation(f, func(t *testing.T, seed fuzzSeed) {
		t.Run("IsZero + String", func(t *testing.T) {
			// no overflow can occur
			truncateToAvoidOverflow := naturalMaxLen

			str := seed.String(truncateToAvoidOverflow)

			a, err := NewFromString(str)
			require.NoError(t, err)

			sa, err := decimal.NewFromString(str)
			require.NoError(t, err)

			require.Equal(t, a.IsZero(), sa.IsZero())
			require.Equal(t, a.String(), sa.String())
		})
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
