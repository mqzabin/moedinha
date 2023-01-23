package currency

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCurrency(t *testing.T) {
	as := strings.Repeat("7", natMaxDigitsPerInt-1)
	as += strings.Repeat("9", natMaxDigitsPerInt)
	as += strings.Repeat("9", natMaxDigitsPerInt-2)
	as += "80"
	as = as[:len(as)-currDecimalDigits] + "." + as[len(as)-currDecimalDigits:]

	bs := strings.Repeat("0", natMaxDigitsPerInt*natNumberOfInts-2)
	bs += "20"
	bs = bs[:len(bs)-currDecimalDigits] + "." + bs[len(bs)-currDecimalDigits:]

	a, err := FromDecimalString(as)
	if err != nil {
		panic(err)
	}

	b, err := FromDecimalString(bs)
	if err != nil {
		panic(err)
	}

	t.Log(a.String() + " +")
	t.Log(b.String() + " =")
	t.Log(a.Add(b).String())
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

	var (
		mCurrency Currency
		sCurrency decimal.Decimal
	)

	b.Run("Add - PoC", func(b *testing.B) {
		x, _ := FromDecimalString(as)

		y, _ := FromDecimalString(bs)

		for i := 0; i < b.N; i++ {

			mCurrency = x.Add(y)
		}
	})

	b.Run("Add - shopspring", func(b *testing.B) {
		x, _ := decimal.NewFromString(as)

		y, _ := decimal.NewFromString(bs)

		for i := 0; i < b.N; i++ {

			sCurrency = x.Add(y)
		}
	})

	b.Run("GreaterThan - PoC", func(b *testing.B) {
		x, _ := FromDecimalString(as)

		y, _ := FromDecimalString(bs)

		for i := 0; i < b.N; i++ {
			if !x.GreaterThan(y) {
				break
			}
		}
	})

	b.Run("GreaterThan - shopspring", func(b *testing.B) {
		x, _ := decimal.NewFromString(as)
		y, _ := decimal.NewFromString(bs)

		for i := 0; i < b.N; i++ {

			if !x.GreaterThan(y) {
				break
			}
		}
	})

	b.Log(mCurrency.String())
	b.Log(sCurrency.String())
}
