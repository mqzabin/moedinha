package moedinha

const zeroRune = '0'

var zeroFiller = newZeroFiller()

func newZeroFiller() [currencyMaxLen]byte {
	var zf [currencyMaxLen]byte
	for i := 0; i < len(zf); i++ {
		zf[i] = zeroRune
	}

	return zf
}
