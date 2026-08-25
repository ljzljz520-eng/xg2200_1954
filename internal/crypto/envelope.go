package crypto

import (
	"fmt"
	"strings"
)

type Envelope struct {
	Algorithm string `json:"algorithm"`
	Digest    string `json:"digest"`
	Body      string `json:"body"`
}

func (c Cipher) Wrap(payload string) Envelope {
	return Envelope{Algorithm: "xor-sha256", Digest: Digest(payload), Body: c.Seal(payload)}
}

func (c Cipher) Unwrap(envelope Envelope) (string, error) {
	payload, err := c.Open(envelope.Body)
	if err != nil {
		return "", err
	}
	if Digest(payload) != envelope.Digest {
		return "", fmt.Errorf("payload digest mismatch")
	}
	return payload, nil
}

func ParseEnvelope(value string) (Envelope, error) {
	parts := strings.SplitN(value, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return Envelope{}, fmt.Errorf("invalid envelope")
	}
	return Envelope{Algorithm: parts[0], Digest: parts[1], Body: parts[2]}, nil
}

func FormatEnvelope(e Envelope) string {
	return strings.Join([]string{e.Algorithm, e.Digest, e.Body}, ":")
}
