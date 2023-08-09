// Package pbkdf2 implements a modular crypt format for PBKDF2-SHA1,
// PBKDF2-SHA256 and PBKDF-SHA512.
//
// The format is the same as that used by Python's passlib and is compatible.
package pbkdf2

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"gopkg.in/hlandau/passlib.v1/abstract"
	"gopkg.in/hlandau/passlib.v1/hash/pbkdf2/raw"
	"hash"
	"strings"
)

// An implementation of Scheme implementing a number of PBKDF2 modular crypt
// formats used by Python's passlib ($pbkdf2$, $pbkdf2-sha256$,
// $pbkdf2-sha512$).
//
// Uses RecommendedRounds.
//
// WARNING: SHA1 should not be used for new applications under any
// circumstances. It should be used for legacy compatibility only.
var SHA1Crypter abstract.Scheme
var SHA256Crypter abstract.Scheme
var SHA512Crypter abstract.Scheme

const (
	RecommendedRoundsSHA1   = 131000
	RecommendedRoundsSHA256 = 29000
	RecommendedRoundsSHA512 = 25000
)

const SaltLength = 16

func init() {
	SHA1Crypter = New("$pbkdf2$", sha1.New, RecommendedRoundsSHA1)
	SHA256Crypter = New("$pbkdf2-sha256$", sha256.New, RecommendedRoundsSHA256)
	SHA512Crypter = New("$pbkdf2-sha512$", sha512.New, RecommendedRoundsSHA512)
}

type scheme struct {
	Ident    string
	HashFunc func() hash.Hash
	Rounds   int
}

func New(ident string, hf func() hash.Hash, rounds int) abstract.Scheme {
	return &scheme{
		Ident:    ident,
		HashFunc: hf,
		Rounds:   rounds,
	}
}

func (s *scheme) Hash(password string) (string, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := raw.Hash([]byte(password), salt, s.Rounds, s.HashFunc)

	newHash := fmt.Sprintf("%s%d$%s$%s", s.Ident, s.Rounds, raw.Base64Encode(salt), hash)
	return newHash, nil
}

func (s *scheme) Verify(password, stub string) (err error) {
	_, rounds, salt, oldHash, err := raw.Parse(stub)
	if err != nil {
		return
	}

	newHash := raw.Hash([]byte(password), salt, rounds, s.HashFunc)

	if len(newHash) == 0 || !abstract.SecureCompare(oldHash, newHash) {
		err = abstract.ErrInvalidPassword
	}

	return
}

func (s *scheme) SupportsStub(stub string) bool {
	return strings.HasPrefix(stub, s.Ident)
}

func (s *scheme) NeedsUpdate(stub string) bool {
	_, rounds, salt, _, err := raw.Parse(stub)
	return err == raw.ErrInvalidRounds || rounds < s.Rounds || len(salt) < SaltLength
}
