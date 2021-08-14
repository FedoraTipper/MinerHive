package hash

import "crypto/md5"

type Md5Hasher struct{}

func (m *Md5Hasher) Hash(data ...[]byte) []byte {
	hash := md5.New()

	for _, d := range data {
		hash.Write(d)
	}

	return hash.Sum(nil)
}
