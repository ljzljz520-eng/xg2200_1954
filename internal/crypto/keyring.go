package crypto

import (
	"fmt"
	"sort"
)

type Keyring struct {
	keys   map[string]Cipher
	active string
}

func NewKeyring() *Keyring {
	return &Keyring{keys: make(map[string]Cipher)}
}

func (k *Keyring) Add(id, secret string) error {
	cipher, err := New(secret)
	if err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("key id is required")
	}
	k.keys[id] = cipher
	if k.active == "" {
		k.active = id
	}
	return nil
}

func (k *Keyring) Activate(id string) error {
	if _, ok := k.keys[id]; !ok {
		return fmt.Errorf("key %s is not registered", id)
	}
	k.active = id
	return nil
}

func (k *Keyring) Seal(payload string) (string, error) {
	cipher, ok := k.keys[k.active]
	if !ok {
		return "", fmt.Errorf("no active key")
	}
	return k.active + "." + cipher.Seal(payload), nil
}

func (k *Keyring) Open(value string) (string, error) {
	parts := splitKey(value)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid keyed payload")
	}
	cipher, ok := k.keys[parts[0]]
	if !ok {
		return "", fmt.Errorf("key %s is not registered", parts[0])
	}
	return cipher.Open(parts[1])
}

func (k *Keyring) IDs() []string {
	ids := make([]string, 0, len(k.keys))
	for id := range k.keys {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func splitKey(value string) []string {
	for index, char := range value {
		if char == '.' {
			return []string{value[:index], value[index+1:]}
		}
	}
	return nil
}
