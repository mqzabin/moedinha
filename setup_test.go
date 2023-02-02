package moedinha

import (
	"reflect"
	"strings"
	"testing"
)

type fuzzSeed struct {
	n   [numberOfUints]uint64
	neg bool
}

// fuzzAddArgs return the args that should be used in f.Add() to represent the seed.
func (fs fuzzSeed) fuzzAddArgs() []any {
	args := make([]any, 0, numberOfUints+1)

	for j := range fs.n {
		args = append(args, fs.n[j])
	}

	args = append(args, fs.neg)

	return args
}

func (fs fuzzSeed) string(truncateTo int) string {
	var result string

	for i := range fs.n {
		fs.n[i], _ = rebalance(fs.n[i], 0)
		nStr := itoa(fs.n[i])
		result += string(nStr[:])
	}

	result = result[:truncateTo]
	if truncateTo < currencyDecimalDigits {
		result = strings.Repeat(string(zeroRune), currencyDecimalDigits-truncateTo) + result
		truncateTo = currencyDecimalDigits
	}

	result = result[:truncateTo-currencyDecimalDigits] + string(currencyDecimalSeparatorSymbol) + result[truncateTo-currencyDecimalDigits:]

	result = strings.TrimLeft(result, string(zeroRune))

	if result[0] == currencyDecimalSeparatorSymbol {
		result = string(zeroRune) + result
	}

	if fs.neg {
		result = "-" + result
	}

	return result
}

func generateMaxArray() [numberOfUints]uint64 {
	var seed [numberOfUints]uint64
	for i := range seed {
		seed[i] = maxValuePerUint
	}

	return seed
}

func generateSeeds() []fuzzSeed {
	seeds := make([]fuzzSeed, 0, 4*numberOfUints)

	// Upper triangular matrix
	for i := 0; i < numberOfUints; i++ {
		var seedN [numberOfUints]uint64

		for j := i; j < numberOfUints; j++ {
			seedN[j] = maxValuePerUint
		}

		seeds = append(seeds, fuzzSeed{n: seedN, neg: false})
	}

	// Lower triangular matrix
	for i := 0; i < numberOfUints; i++ {
		seedN := generateMaxArray()

		for j := i; j < numberOfUints; j++ {
			seedN[j] = 0
		}

		seeds = append(seeds, fuzzSeed{n: seedN, neg: false})
	}

	// Doubling for negative values
	for i := range seeds {
		seeds = append(seeds, fuzzSeed{n: seeds[i].n, neg: true})
	}

	return seeds
}

func fuzzyUnaryOperation(f *testing.F, fn func(*testing.T, fuzzSeed)) {
	f.Helper()

	seeds := generateSeeds()

	for _, seed := range seeds {
		f.Add(seed.fuzzAddArgs()...)
	}

	typeUint64 := reflect.TypeOf(uint64(1))
	typeBool := reflect.TypeOf(true)

	seedArgs := make([]reflect.Type, 0, numberOfUints+1)
	for i := 0; i < numberOfUints; i++ {
		seedArgs = append(seedArgs, typeUint64)
	}
	seedArgs = append(seedArgs, typeBool)

	funcParameters := append([]reflect.Type{reflect.TypeOf(&testing.T{})}, seedArgs...)

	funcSignature := reflect.FuncOf(funcParameters, []reflect.Type{}, false)

	fuzzFunc := reflect.MakeFunc(funcSignature, func(args []reflect.Value) (results []reflect.Value) {
		t := args[0].Interface().(*testing.T)
		t.Parallel()

		var a fuzzSeed
		for i := range a.n {
			a.n[i] = args[i+1].Interface().(uint64)
		}
		a.neg = args[numberOfUints+1].Interface().(bool)

		fn(t, a)

		return nil
	})

	f.Fuzz(fuzzFunc.Interface())
}

func fuzzyBinaryOperation(f *testing.F, fn func(*testing.T, fuzzSeed, fuzzSeed)) {
	f.Helper()

	seeds := generateSeeds()

	for _, seedA := range seeds {
		for _, seedB := range seeds {
			f.Add(append(
				seedA.fuzzAddArgs(),
				seedB.fuzzAddArgs()...,
			)...)
		}
	}

	typeUint64 := reflect.TypeOf(uint64(1))
	typeBool := reflect.TypeOf(true)

	seedArgs := make([]reflect.Type, 0, numberOfUints+1)
	for i := 0; i < numberOfUints; i++ {
		seedArgs = append(seedArgs, typeUint64)
	}
	seedArgs = append(seedArgs, typeBool)

	funcParameters := append(append([]reflect.Type{reflect.TypeOf(&testing.T{})}, seedArgs...), seedArgs...)

	funcSignature := reflect.FuncOf(funcParameters, []reflect.Type{}, false)

	fuzzFunc := reflect.MakeFunc(funcSignature, func(args []reflect.Value) (results []reflect.Value) {
		t := args[0].Interface().(*testing.T)
		t.Parallel()

		var a fuzzSeed
		for i := range a.n {
			a.n[i] = args[i+1].Interface().(uint64)
		}
		a.neg = args[numberOfUints+1].Interface().(bool)

		var b fuzzSeed
		for i := range b.n {
			b.n[i] = args[i+numberOfUints+2].Interface().(uint64)
		}
		b.neg = args[2*(numberOfUints+1)].Interface().(bool)

		fn(t, a, b)

		return nil
	})

	f.Fuzz(fuzzFunc.Interface())
}
