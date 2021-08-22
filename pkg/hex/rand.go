package hex

import "math/rand"

const (
	HexValues = "0123456789ABCDEF"
	HexLength = 16
)

func GenerateHexString(length int) string {
	var hex []rune
	for i := 0; i < length; i++ {
		idx := rand.Intn(HexLength)
		hex = append(hex, rune(HexValues[idx]))
	}

	return string(hex)
}
