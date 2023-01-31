# moedinha

```
goos: linux
goarch: amd64
pkg: github.com/mqzabin/moedinha
cpu: AMD Ryzen 7 5700G with Radeon Graphics         
BenchmarkNewFromString/moedinha-16               1803176               670.0 ns/op             0 B/op          0 allocs/op
BenchmarkNewFromString/shopspring-16             1421364               737.5 ns/op           184 B/op          5 allocs/op
BenchmarkString/moedinha-16                      6693012               173.8 ns/op            64 B/op          1 allocs/op
BenchmarkString/shopspring-16                    3155907               441.3 ns/op           320 B/op          5 allocs/op
BenchmarkAdd/moedinha-16                        352125472                3.373 ns/op           0 B/op          0 allocs/op
BenchmarkAdd/shopspring-16                       3686824               369.2 ns/op           304 B/op          8 allocs/op
BenchmarkSub/moedinha-16                        215181790                5.567 ns/op           0 B/op          0 allocs/op
BenchmarkSub/shopspring-16                       2553356               436.4 ns/op           304 B/op          8 allocs/op
PASS
ok      github.com/mqzabin/moedinha     13.502s
```