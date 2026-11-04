package identifier

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	Prefix   = "RP-"
	Length   = 8
	Alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func Generate() (string, error) {
	id, err := gonanoid.Generate(Alphabet, Length)
	if err != nil {
		return "", err
	}
	return Prefix + id, nil
}
