package hex

import "math/rand"

const (
	possibleHexValues = "0123456789ABCDEF"
)

func GenerateHexString(length int) string {
	var hex []rune
	for i := 0; i < length; i++ {
		idx := rand.Int()
		hex = append(hex, rune(possibleHexValues[idx]))
	}

	return string(hex)
}
