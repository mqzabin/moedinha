FUZZ_PARALLELISM=8
FUZZ_CACHE_DIR=testdata

.PHONY: fuzz/unary
fuzz/unary:
	@go test -fuzz=FuzzUnary -parallel=$(FUZZ_PARALLELISM) -test.fuzzcachedir=$(FUZZ_CACHE_DIR)

.PHONY: fuzz/comparisons
fuzz/comparisons:
	@go test -fuzz=FuzzComparisons -parallel=$(FUZZ_PARALLELISM) -test.fuzzcachedir=$(FUZZ_CACHE_DIR)

.PHONY: fuzz/addsub
fuzz/addsub:
	@go test -fuzz=FuzzAddSub -parallel=$(FUZZ_PARALLELISM) -test.fuzzcachedir=$(FUZZ_CACHE_DIR)

.PHONY: fuzz/mul
fuzz/mul:
	@go test -fuzz=FuzzMul -parallel=$(FUZZ_PARALLELISM) -test.fuzzcachedir=$(FUZZ_CACHE_DIR)

.PHONY: fuzz/clean
fuzz/clean:
	@go clean -fuzzcache

.PHONY: bench
bench:
	@go test -run none -bench=. -benchmem ./...
