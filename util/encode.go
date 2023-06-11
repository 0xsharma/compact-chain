package util

import (
	"bytes"
	"encoding/gob"
	"log"
)

func EncodeToBytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func DecodeFromBytes[T any](in []byte) (*T, error) {
	out := new(T)
	dec := gob.NewDecoder(bytes.NewReader(in))
	err := dec.Decode(&out)
	return out, err
}
