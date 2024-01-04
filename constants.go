package moedinha

const (
	highLow = 2
	low     = 1
	high    = 0
)

var (
	// Natural constants
	naturalOne = naturalFromUint(1)
)

func naturalFromUint(v uint64) natural {
	n := natural{}
	n[numberOfUints-1] = v

	return n
}
