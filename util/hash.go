package util

const (
	hashLength    = 32
	addressLength = 20
)

type Hash [hashLength]byte
type Address [addressLength]byte

func (h Hash) String() string {
	return string(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (a Address) String() string {
	return string(a[:])
}

func (a Address) Bytes() []byte {
	return a[:]
}
