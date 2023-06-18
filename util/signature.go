package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

type UnlockedAccount struct {
	privateKey *ecdsa.PrivateKey
}

func NewUnlockedAccount(privateKey *ecdsa.PrivateKey) *UnlockedAccount {
	return &UnlockedAccount{
		privateKey: privateKey,
	}
}

func (ua *UnlockedAccount) PublicKey() *ecdsa.PublicKey {
	return &ua.privateKey.PublicKey
}

func PublicKeyToAddress(pubkey *ecdsa.PublicKey) *Address {
	data := elliptic.Marshal(pubkey, pubkey.X, pubkey.Y)
	return NewAddress(data)
}

func (ua *UnlockedAccount) Address() *Address {
	pubkey := ua.privateKey.PublicKey
	return PublicKeyToAddress(&pubkey)
}

func (ua *UnlockedAccount) Sign(data []byte) (*big.Int, *big.Int, error) {
	return ecdsa.Sign(rand.Reader, ua.privateKey, data)
}

func (ua *UnlockedAccount) Verify(data []byte, r *big.Int, s *big.Int) bool {
	return ecdsa.Verify(ua.PublicKey(), data, r, s)
}

func HexToPrivateKey(hexStr string) *ecdsa.PrivateKey {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}

	k := new(big.Int)
	k.SetBytes(bytes)

	priv := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	priv.PublicKey.Curve = curve
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())

	return priv
}
