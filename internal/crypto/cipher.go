package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Cipher struct {
	key []byte
}

func New(key string) (Cipher, error) {
	if len(key) < 8 {
		return Cipher{}, errors.New("key must contain at least eight characters")
	}
	return Cipher{key: []byte(key)}, nil
}

func (c Cipher) Seal(plain string) string {
	out := make([]byte, len(plain))
	for i := range []byte(plain) {
		out[i] = plain[i] ^ c.key[i%len(c.key)]
	}
	return hex.EncodeToString(out)
}

func (c Cipher) Open(encoded string) (string, error) {
	data, err := hex.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	plain := make([]byte, len(data))
	for i := range data {
		plain[i] = data[i] ^ c.key[i%len(c.key)]
	}
	return string(plain), nil
}

func Digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func Fingerprint(values ...string) string {
	h := sha256.New()
	for _, value := range values {
		h.Write([]byte(value))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}
