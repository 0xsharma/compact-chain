package util

import (
	"crypto/sha256"
	"fmt"
)

const (
	hashLength    = 32
	addressLength = 20
)

type Hash [hashLength]byte
type Address [addressLength]byte

func (h Hash) String() string {
	hexString := fmt.Sprintf("0x%x", h[:])
	return hexString
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (a Address) String() string {
	hexString := fmt.Sprintf("0x%x", a[:])
	return hexString
}

func (a Address) Bytes() []byte {
	return a[:]
}

func ByteToHash(b []byte) *Hash {
	hash := Hash{}
	copy(hash[:], b[:])

	return &hash
}

func NewHash(b []byte) *Hash {
	sum := sha256.Sum256(b)

	hash := Hash{}
	copy(hash[:], sum[:])

	return &hash
}

func NewHashFromHex(s string) *Hash {
	hash := NewHash([]byte(s))

	return hash
}

func NewAddress(b []byte) *Address {
	sum := sha256.Sum256(b)

	address := Address{}
	copy(address[:], sum[:])

	return &address
}

func NewAddressFromHex(s string) *Address {
	address := NewAddress([]byte(s))

	return address
}
