package raw

import (
	"golang.org/x/crypto/pbkdf2"
	"hash"
)

const (
	MinRounds = 1
	MaxRounds = 0x7fffffff // setting at 32-bit signed integer limit for now
)

func Hash(password, salt []byte, rounds int, hf func() hash.Hash) (hash string) {
	return Base64Encode(pbkdf2.Key(password, salt, rounds, hf().Size(), hf))
}
