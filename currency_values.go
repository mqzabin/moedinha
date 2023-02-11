package moedinha

var (
	currencyOne     = currencyFromUint(1)
	currencyTwo     = currencyFromUint(2)
	currencyOneHalf = aHalf()
)

func currencyFromUint(v uint64) Currency {
	c := Currency{}
	c.t.n[numberOfUints-uintsReservedToDecimal-1] = v

	return c
}

func aHalf() Currency {
	c := Currency{}
	c.t.n[numberOfUints-uintsReservedToDecimal] = 500000000000000000

	return c
}
