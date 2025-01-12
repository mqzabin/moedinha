package moedinha

import (
	"strconv"
	"testing"

	"github.com/mqzabin/fuzzdecimal"
	"github.com/shopspring/decimal"
)

func FuzzUnary(f *testing.F) {
	parseDecimal := func(t *fuzzdecimal.T, s string) (Currency, error) {
		t.Helper()

		return NewFromString(s)
	}

	parseShopspringDecimal := func(t *fuzzdecimal.T, s string) (decimal.Decimal, error) {
		t.Helper()

		return decimal.NewFromString(s)
	}

	fuzzdecimal.Fuzz(f, 1, func(t *fuzzdecimal.T) {
		fuzzdecimal.AsDecimalComparison1(t, "IsZero", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1 decimal.Decimal) (string, error) {
				t.Helper()

				return strconv.FormatBool(x1.IsZero()), nil
			},
			func(t *fuzzdecimal.T, x1 Currency) string {
				return strconv.FormatBool(x1.IsZero())
			},
		)

		fuzzdecimal.AsDecimalComparison1(t, "String", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1 decimal.Decimal) (string, error) {
				t.Helper()

				return x1.String(), nil
			},
			func(t *fuzzdecimal.T, x1 Currency) string {
				return x1.String()
			},
		)
	}, fuzzdecimal.WithAllDecimals(
		fuzzdecimal.WithSigned(),
		fuzzdecimal.WithMaxSignificantDigits(naturalMaxLen),
		fuzzdecimal.WithDecimalPointAt(currencyDecimalDigits),
	))
}

func FuzzComparisons(f *testing.F) {
	parseDecimal := func(t *fuzzdecimal.T, s string) (Currency, error) {
		t.Helper()

		return NewFromString(s)
	}

	parseShopspringDecimal := func(t *fuzzdecimal.T, s string) (decimal.Decimal, error) {
		t.Helper()

		return decimal.NewFromString(s)
	}

	fuzzdecimal.Fuzz(f, 2, func(t *fuzzdecimal.T) {
		fuzzdecimal.AsDecimalComparison2(t, "Equal", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return strconv.FormatBool(x1.Equal(x2)), nil
			},
			func(t *fuzzdecimal.T, x1, x2 Currency) string {
				return strconv.FormatBool(x1.Equal(x2))
			},
		)

		fuzzdecimal.AsDecimalComparison2(t, "GreaterThan", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return strconv.FormatBool(x1.GreaterThan(x2)), nil
			},
			func(t *fuzzdecimal.T, x1, x2 Currency) string {
				return strconv.FormatBool(x1.GreaterThan(x2))
			},
		)

		fuzzdecimal.AsDecimalComparison2(t, "GreaterThanOrEqual", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return strconv.FormatBool(x1.GreaterThanOrEqual(x2)), nil
			},
			func(t *fuzzdecimal.T, x1, x2 Currency) string {
				return strconv.FormatBool(x1.GreaterThanOrEqual(x2))
			},
		)

		fuzzdecimal.AsDecimalComparison2(t, "LessThan", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return strconv.FormatBool(x1.LessThan(x2)), nil
			},
			func(t *fuzzdecimal.T, x1, x2 Currency) string {
				return strconv.FormatBool(x1.LessThan(x2))
			},
		)

		fuzzdecimal.AsDecimalComparison2(t, "LessThanOrEqual", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return strconv.FormatBool(x1.LessThanOrEqual(x2)), nil
			},
			func(t *fuzzdecimal.T, x1, x2 Currency) string {
				return strconv.FormatBool(x1.LessThanOrEqual(x2))
			},
		)
	}, fuzzdecimal.WithAllDecimals(
		fuzzdecimal.WithSigned(),
		fuzzdecimal.WithMaxSignificantDigits(naturalMaxLen),
		fuzzdecimal.WithDecimalPointAt(currencyDecimalDigits),
	))
}

func FuzzAddSub(f *testing.F) {
	parseDecimal := func(t *fuzzdecimal.T, s string) (Currency, error) {
		t.Helper()

		return NewFromString(s)
	}

	parseShopspringDecimal := func(t *fuzzdecimal.T, s string) (decimal.Decimal, error) {
		t.Helper()

		return decimal.NewFromString(s)
	}

	fuzzdecimal.Fuzz(f, 2, func(t *fuzzdecimal.T) {
		fuzzdecimal.AsDecimalComparison2(t, "Add", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return x1.Add(x2).Truncate(currencyDecimalDigits).String(), nil
			},
			func(t *fuzzdecimal.T, x1, x2 Currency) string {
				return x1.Add(x2).String()
			},
		)

		fuzzdecimal.AsDecimalComparison2(t, "Sub", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return x1.Sub(x2).Truncate(currencyDecimalDigits).String(), nil
			},
			func(t *fuzzdecimal.T, x1 Currency, x2 Currency) string {
				return x1.Sub(x2).String()
			},
		)
	}, fuzzdecimal.WithAllDecimals(
		fuzzdecimal.WithSigned(),
		// The a+b and a-b operations will at most add 1 digit to the greatest number between a and b.
		// So, we should ensure that the greatest number has at most naturalMaxLen-1 digits.
		fuzzdecimal.WithMaxSignificantDigits(naturalMaxLen-1),
		fuzzdecimal.WithDecimalPointAt(currencyDecimalDigits),
	))
}

func FuzzMul(f *testing.F) {
	parseDecimal := func(t *fuzzdecimal.T, s string) (Currency, error) {
		t.Helper()

		return NewFromString(s)
	}

	parseShopspringDecimal := func(t *fuzzdecimal.T, s string) (decimal.Decimal, error) {
		t.Helper()

		return decimal.NewFromString(s)
	}

	fuzzdecimal.Fuzz(f, 2, func(t *fuzzdecimal.T) {
		fuzzdecimal.AsDecimalComparison2(t, "Add", parseDecimal, parseShopspringDecimal,
			func(t *fuzzdecimal.T, x1, x2 decimal.Decimal) (string, error) {
				t.Helper()

				return x1.Mul(x2).Truncate(currencyDecimalDigits).String(), nil
			},
			func(t *fuzzdecimal.T, x1 Currency, x2 Currency) string {
				return x1.Mul(x2).String()
			},
		)
	}, fuzzdecimal.WithAllDecimals(
		fuzzdecimal.WithSigned(),
		// Multiplication result will at most sum the number of digits of "a" and "b" in "a*b".
		// So, we should ensure that digits(a) + digits(b) don't overflow the
		// naturalMaxLen constant.
		fuzzdecimal.WithMaxSignificantDigits(naturalMaxLen/2),
		fuzzdecimal.WithDecimalPointAt(currencyDecimalDigits),
	))
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

func BenchmarkMul(b *testing.B) {
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
