package hash

type Hasher interface {
	Hash(data ...[]byte) []byte
}

const (
	Keccak = "keccak"
	Md5    = "md5"
	Sha256 = "sha256"
)

func GetHasher(hasher string) Hasher {
	switch hasher {
	case Keccak:
		return &Keccak256Hasher{}
	case Md5:
		return &Md5Hasher{}
	case Sha256:
		return &Sha256Hasher{}
	default:
		panic("Failed to find relevant hasher")
	}
}
