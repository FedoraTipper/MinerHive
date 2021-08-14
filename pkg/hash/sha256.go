package hash

import (
	"crypto/sha256"
)

type Sha256Hasher struct{}

func (m *Sha256Hasher) Hash(data ...[]byte) []byte {
	hash := sha256.New()

	for _, d := range data {
		hash.Write(d)
	}

	return hash.Sum(nil)
}
