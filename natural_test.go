package moedinha

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_natural_divByUint(t *testing.T) {

	createNatural := func(s string) natural {

		var strInput [naturalMaxLen]byte
		copy(strInput[naturalMaxLen-len(s):], s)
		copy(strInput[:naturalMaxLen-len(s)], zeroFiller[:naturalMaxLen-len(s)])

		n, err := newNatFromString(strInput)
		require.NoError(t, err)

		return n
	}

	stringNatural := func(n natural) string {
		if n.isZero() {
			return "0"
		}

		gotString1 := n.string()
		return strings.TrimLeft(string(gotString1[:]), "0")
	}

	t.Run("div by uint", func(t *testing.T) {
		tests := []struct {
			name string
			n    string
			d    uint64
		}{
			{
				name: "",
				n:    "999999999999999999999999999999999999999999999999999999999999999999999999",
				d:    999999999999999999,
			},
			{
				name: "",
				n:    "000000000000000000000000000000000000099999999999999999999999999999999999",
				d:    999999999999999999,
			},
			{
				name: "",
				n:    "000000000000000000000000000000000000099999999999999999999999999999999999",
				d:    3,
			},
			{
				name: "",
				n:    "100000000000000000000000000000000000000000000000000000000000000000000000",
				d:    6,
			},
			{
				name: "",
				n:    "100000000000000000",
				d:    6,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var (
					wantDiv string
					wantRem uint64
				)
				{
					a, err := decimal.NewFromString(tt.n)
					require.NoError(t, err)

					b := decimal.NewFromInt(int64(tt.d))

					c := a.Div(b).Truncate(0)

					wantDiv = c.String()

					wantRem = uint64(a.Sub(b.Mul(c)).Truncate(0).IntPart())
				}

				n := createNatural(tt.n)

				got, got1 := n.divByUint(tt.d)

				assert.Equal(t, wantDiv, stringNatural(got))
				assert.Equal(t, wantRem, got1)
			})
		}
	})

	t.Run("div by nat", func(t *testing.T) {
		tests := []struct {
			name string
			n    string
			d    string
		}{
			{
				name: "",
				n:    "89999",
				d:    "990",
			},
			{
				name: "",
				n:    "999999999999999999999999999999999999",
				d:    "999999999999999999",
			},
			{
				name: "",
				n:    "999999999999999999999999999999999999999999999999999999999999999999999999",
				d:    "999999999999999999",
			},
			{
				name: "",
				n:    "1234234213423545367456123123",
				d:    "6",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var wantDiv, wantRem string
				{
					a, err := decimal.NewFromString(tt.n)
					require.NoError(t, err)

					b, err := decimal.NewFromString(tt.d)
					require.NoError(t, err)

					c := a.Div(b).Truncate(0)

					wantDiv = c.String()

					wantRem = a.Sub(b.Mul(c)).Truncate(0).String()
				}

				n := createNatural(tt.n)

				d := createNatural(tt.d)

				fmt.Println("n / d = ", stringNatural(n), "/", stringNatural(d))

				got, got1 := n.div(d)

				fmt.Println(stringNatural(got), "rem", stringNatural(got1))

				assert.Equal(t, wantDiv, stringNatural(got))
				assert.Equal(t, wantRem, stringNatural(got1))
			})
		}
	})

}
