package moedinha

import (
	"testing"
)

func Test_multiplyUint(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name     string
		args     args
		digits   uint64
		overflow uint64
	}{{
		name: "lalala",
		args: args{
			a: maxValuePerUint,
			b: maxValuePerUint,
		},
		digits:   000000000000000001,
		overflow: 999999999999999998,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digits, overflow := multiplyUint(tt.args.a, tt.args.b)
			t.Log(digits)
			t.Log(overflow)
			if digits != tt.digits {
				t.Errorf("multiplyUint() digits = %v, want %v", digits, tt.digits)
			}
			if overflow != tt.overflow {
				t.Errorf("multiplyUint() overflow = %v, want %v", overflow, tt.overflow)
			}
		})
	}
}
