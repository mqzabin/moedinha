package moedinha

import (
	"testing"
)

func Test_highLowDiv(t *testing.T) {
	type args struct {
		highDigits  uint64
		lowDigits   uint64
		denominator uint64
	}
	tests := []struct {
		name    string
		args    args
		wantQuo uint64
		wantRem uint64
	}{
		{
			name: "",
			args: args{
				highDigits:  99999999999999999,
				lowDigits:   999999999999999999,
				denominator: 999999999999999999,
			},
			wantQuo: 100000000000000000,
			wantRem: 99999999999999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuo, gotRem := highLowDiv(tt.args.highDigits, tt.args.lowDigits, tt.args.denominator)
			if gotQuo != tt.wantQuo {
				t.Errorf("highLowDiv() gotQuo = %v, wantQuo %v", gotQuo, tt.wantQuo)
			}
			if gotRem != tt.wantRem {
				t.Errorf("highLowDiv() gotRem = %v, wantQuo %v", gotRem, tt.wantRem)
			}
		})
	}
}
