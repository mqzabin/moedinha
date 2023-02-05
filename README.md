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
	a, _ := moedinha.NewFromString("999999999999999999999999999999999999.999999999999999999")
	b, _ := moedinha.NewFromString("888888888888888888.888888888888888888")

	fmt.Printf("a + b = %s", a.Add(b).String())
	fmt.Printf("a - b = %s", a.Sub(b).String())
	fmt.Printf("b - a = %s", b.Sub(a).String())
	fmt.Printf("a * b = %s", a.Mul(b).String())
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

# Roadmap
- [X] Sum.
- [X] Subtraction.
- [X] Multiplication.
- [ ] Division (Newton-Raphson method).
- [ ] Exponentiation.

# Benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/mqzabin/moedinha
cpu: AMD Ryzen 7 5700G with Radeon Graphics         
BenchmarkNewFromString/moedinha-16               1643692               735.1 ns/op             0 B/op          0 allocs/op
BenchmarkNewFromString/shopspring-16             1362316               844.2 ns/op           184 B/op          5 allocs/op
BenchmarkString/moedinha-16                      4755722               239.9 ns/op            64 B/op          1 allocs/op
BenchmarkString/shopspring-16                    2683314               482.7 ns/op           320 B/op          5 allocs/op
BenchmarkAdd/moedinha-16                        29274129                39.17 ns/op            0 B/op          0 allocs/op
BenchmarkAdd/shopspring-16                       2546943               426.8 ns/op           304 B/op          8 allocs/op
BenchmarkSub/moedinha-16                         3190180               371.5 ns/op             0 B/op          0 allocs/op
BenchmarkSub/shopspring-16                      13402520               132.8 ns/op            96 B/op          2 allocs/op
BenchmarkMul/moedinha-16                        17169872                63.04 ns/op            0 B/op          0 allocs/op
BenchmarkMul/shopspring-16                       3261558               438.9 ns/op           304 B/op          8 allocs/op

```