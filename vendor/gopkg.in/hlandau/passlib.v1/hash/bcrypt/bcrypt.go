// Package bcrypt implements the bcrypt password hashing mechanism.
//
// Please note that bcrypt truncates passwords to 72 characters in length. Consider using
// a more modern hashing scheme such as scrypt or sha-crypt. If you must use bcrypt,
// consider using bcrypt-sha256 instead.
package bcrypt

import "golang.org/x/crypto/bcrypt"
import "gopkg.in/hlandau/passlib.v1/abstract"
import "fmt"

// An implementation of Scheme implementing bcrypt.
//
// Uses RecommendedCost.
var Crypter abstract.Scheme

// The recommended cost for bcrypt. This may change with subsequent releases.
const RecommendedCost = 12

// bcrypt.DefaultCost is a bit low (10), so use 12 instead.

func init() {
	Crypter = New(RecommendedCost)
}

// Create a new scheme implementing bcrypt. The recommended cost is RecommendedCost.
func New(cost int) abstract.Scheme {
	return &scheme{
		Cost: cost,
	}
}

type scheme struct {
	Cost int
}

func (s *scheme) SupportsStub(stub string) bool {
	return len(stub) >= 3 && stub[0] == '$' && stub[1] == '2' &&
		(stub[2] == '$' || (len(stub) >= 4 && stub[3] == '$' &&
			(stub[2] == 'a' || stub[2] == 'b' || stub[2] == 'y')))
}

func (s *scheme) Hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), s.Cost)
	if err != nil {
		return "", err
	}

	return string(h), nil
}

func (s *scheme) Verify(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		err = abstract.ErrInvalidPassword
	}

	return err
}

func (s *scheme) NeedsUpdate(stub string) bool {
	cost, err := bcrypt.Cost([]byte(stub))
	if err != nil {
		return false
	}

	return cost < s.Cost
}

func (s *scheme) String() string {
	return fmt.Sprintf("bcrypt(%d)", s.Cost)
}
