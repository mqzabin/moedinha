# moedinha

[![Go Report Card](https://goreportcard.com/badge/github.com/mqzabin/moedinha)](https://goreportcard.com/report/github.com/mqzabin/moedinha)

Fixed-precision decimal numbers in Go, aiming to represent currency values.

# Example

```go
package main

import (
	"fmt"

	"github.com/mqzabin/moedinha"
)

func main() {
	a, _ := moedinha.NewFromString("99999999999999999999999999999999999.999999999999999999")
	b, _ := moedinha.NewFromString("888888888888888888.888888888888888888")

	fmt.Println("a + b =", a.Add(b).String())
	// a + b = 100000000000000000888888888888888888.888888888888888887
	fmt.Println("a - b =", a.Sub(b).String())
	// a - b = 99999999999999999111111111111111111.111111111111111111
	fmt.Println("b - a =", b.Sub(a).String())
	// b - a = -99999999999999999111111111111111111.111111111111111111
	fmt.Println("a * b =", a.Mul(b).String())
	// a * b = 88888888888888888888888888888888888799999999999999999.111111111111111111
}
```

# How it works?

`moedinha` uses an array of `uint64` to represent decimal numbers. Each `uint64` represents up to 18-digits.

You can set how many `uint64` you want to use, and how many of those you want to use as decimal digits. Those settings are
set through editing the [settings.go](./settings.go) file.

The default setting is to use 4 `uint64` and 1 of those as decimal digits. This settings can represent numbers up to:

`999999999999999999999999999999999999999999999999999999.999999999999999999`,

i.e. 54 integer digits and 18 decimal digits.

Since the precision is fixed, overflows during arithmetic operations can happen and the package will call a `panic`.


# Motivation
The [shopspring/decimal](https://github.com/shopspring/decimal) solve the problem of arbitrary precision decimals in Go,
wrapping the `math/big` structure with an easy-to-use API.

The API translation partially loses the `math/big`'s memory allocation optimizations, making many arithmetic operations 
reallocate unnecessary memory. For low throughput scenarios, the garbage collector will handle this and the easy API + 
not worrying about precision/value limits will pay off.

However, in scenarios where there are high throughput and many arithmetical operations per request, the garbage collector
will start to be a performance bottleneck.

These problems arise from the "arbitrary precision" hypothesis, which is not required in some real-world scenarios.
`moedinha` makes this hypothesis false to reduce allocations to zero, using fixed-size arrays to represent integer values.

# Benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/mqzabin/moedinha
cpu: 12th Gen Intel(R) Core(TM) i7-1265U
BenchmarkNewFromString/moedinha-12               1812570               641.2 ns/op             0 B/op          0 allocs/op
BenchmarkNewFromString/shopspring-12             2574252               458.7 ns/op           184 B/op          5 allocs/op
BenchmarkString/moedinha-12                      8340163               149.3 ns/op            64 B/op          1 allocs/op
BenchmarkString/shopspring-12                    4348003               270.1 ns/op           320 B/op          5 allocs/op
BenchmarkAdd/moedinha-12                        34227067                32.50 ns/op            0 B/op          0 allocs/op
BenchmarkAdd/shopspring-12                       5471982               225.6 ns/op           296 B/op          7 allocs/op
BenchmarkSub/moedinha-12                        22223533                54.08 ns/op            0 B/op          0 allocs/op
BenchmarkSub/shopspring-12                       5219170               226.3 ns/op           296 B/op          7 allocs/op
BenchmarkMul/moedinha-12                         4043478               298.6 ns/op             0 B/op          0 allocs/op
BenchmarkMul/shopspring-12                      17446844                60.06 ns/op           96 B/op          2 allocs/op
```

# Development

## Fuzzy tests

This package makes a heavy usage of fuzzy tests, comparing all operation results with their equivalent in the `shopspring/decimal` package.

It uses the [github.com/mqzabin/fuzzdecimal](https://github.com/mqzabin/fuzzdecimal) package to ease the fuzzing process with decimal numbers.

There are some Makefile targets to ease the fuzzy calling process:

- `make fuzz/unary`: Tests unary functions, like `String`, `IsZero`, etc.
- `make fuzz/comparisons`: Tests comparisons functions, like `Equal`, `GreaterThan`, etc.
- `make fuzz/addsub`: Tests `Add` and `Sub` operations.
- `make fuzz/mul`:  Tests `Mul` operations.

All of this target will read and save the fuzzy entries cache to the `./testdata` directory, so the fuzzy process could continue across different machines. 