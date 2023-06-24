package util

import (
	"testing"
)

func TestSigning(t *testing.T) {
	t.Parallel()

	// Address = 0xa52c981eee8687b5e4afd69aa5006548c24d7685
	pkey := HexToPrivateKey("c3fc038a9abc0f483e2e1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a6")
	ua := NewUnlockedAccount(pkey)

	data := []byte("data")

	r, s, err := ua.Sign(data)
	if err != nil {
		t.Fatal(err)
	}

	verified := ua.Verify(data, r, s)

	if !verified {
		t.Fatal("expected true", "got", verified)
	}
}
