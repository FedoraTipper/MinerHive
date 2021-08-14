package hash

import (
	"golang.org/x/crypto/sha3"
)

type Keccak256Hasher struct{}

func (k *Keccak256Hasher) Hash(data ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()

	for _, d := range data {
		hash.Write(d)
	}

	return hash.Sum(nil)
}
