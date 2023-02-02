.PHONY: fuzzy/binary
fuzzy/binary:
	@go test -fuzz=FuzzBinaryOperations -parallel=4

.PHONY: fuzzy/unary
fuzzy/unary:
	@go test -fuzz=FuzzUnaryOperations -parallel=4

.PHONY: fuzzy/clean
fuzzy/clean:
	@go clean -fuzzcache

.PHONY: bench
bench:
	@go test -run none -bench=. -benchmem ./...
