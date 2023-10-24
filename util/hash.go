package util

import (
	"crypto/sha256"
	"fmt"
)

const (
	hashLength    = 32 // Length of hash
	addressLength = 20 // Length of address
)

// Hash is the hash of the block.
type Hash [hashLength]byte

// Address is the address of the account.
type Address [addressLength]byte

// String returns the string representation of the hash.
func (h Hash) String() string {
	hexString := fmt.Sprintf("0x%x", h[:])
	return hexString
}

// Bytes returns the byte representation of the hash.
func (h Hash) Bytes() []byte {
	return h[:]
}

// String returns the string representation of the address.
func (a Address) String() string {
	hexString := fmt.Sprintf("0x%x", a[:])
	return hexString
}

// Bytes returns the byte representation of the address.
func (a Address) Bytes() []byte {
	return a[:]
}

// ByteToHash converts a byte array to hash.
func ByteToHash(b []byte) *Hash {
	hash := Hash{}
	copy(hash[:], b[:])

	return &hash
}

// NewHash creates a new sha256 hash from the given byte array.
func HashData(b []byte) *Hash {
	sum := sha256.Sum256(b)

	hash := Hash{}
	copy(hash[:], sum[:])

	return &hash
}

// NewAddress creates a new address from the given byte array.
func BytesToAddress(b []byte) *Address {
	address := Address{}
	copy(address[:], b[:])

	return &address
}

// string to Address
func StringToAddress(s string) *Address {
	return BytesToAddress([]byte(s))
}
