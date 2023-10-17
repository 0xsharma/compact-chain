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

type CompactPublicKey struct {
	CurveParams *elliptic.CurveParams `json:"Curve"`
	X           *big.Int              `json:"X"`
	Y           *big.Int              `json:"Y"`
}

func (cpk *CompactPublicKey) PublicKey() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: cpk.CurveParams,
		X:     cpk.X,
		Y:     cpk.Y,
	}
}

func PublicKeyToCompact(pubkey *ecdsa.PublicKey) *CompactPublicKey {
	return &CompactPublicKey{
		CurveParams: pubkey.Curve.Params(),
		X:           pubkey.X,
		Y:           pubkey.Y,
	}
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

	address := Address{}
	copy(address[:], data[2:])

	return &address
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
