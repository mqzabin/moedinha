package moedinha

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func Test_highLowDiv(t *testing.T) {
	type args struct {
		highDigits  uint64
		lowDigits   uint64
		denominator uint64
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 uint64
	}{
		{
			name: "",
			args: args{
				highDigits:  99999999999999999,
				lowDigits:   999999999999999999,
				denominator: 999999999999999999,
			},
			want:  1,
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := highLowDiv(tt.args.highDigits, tt.args.lowDigits, tt.args.denominator)
			if got != tt.want {
				t.Errorf("highLowDiv() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("highLowDiv() got1 = %v, want %v", got1, tt.want1)
			}

			a, err := decimal.NewFromString(fmt.Sprintf("%018d%018d", tt.args.highDigits, tt.args.lowDigits))
			if err != nil {
				t.Fail()
			}

			b := decimal.NewFromInt(int64(tt.args.denominator))

			c := a.Div(b).Truncate(0)

			rem := a.Sub(b.Mul(c)).Truncate(0).String()

			d := c.String()

			fmt.Printf("a / b = %s / %s = %s len(%d) rem(%s)\n", a.String(), b.String(), d, len(d), rem)
		})
	}
}
