.PHONY: fuzzy/binary
fuzzy/binary:
	@go test -run=none -fuzz=FuzzBinaryOperations -parallel=4

.PHONY: fuzzy/unary
fuzzy/unary:
	@go test -run=none -fuzz=FuzzUnaryOperations -parallel=4

.PHONY: fuzzy/clean
fuzzy/clean:
	@go clean -fuzzcache

.PHONY: bench
bench:
	@go test -run=none -bench=. -benchmem ./...
