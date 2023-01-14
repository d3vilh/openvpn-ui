// Package bcryptsha256 implements bcrypt with a SHA256 prehash in a format that is compatible with Python passlib's equivalent bcrypt-sha256 scheme.
//
// This is preferred over bcrypt because the prehash essentially renders bcrypt's password length
// limitation irrelevant; although of course it is less compatible.
package bcryptsha256

import "gopkg.in/hlandau/passlib.v1/abstract"
import "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
import "encoding/base64"
import "crypto/sha256"
import "strings"
import "fmt"

type scheme struct {
	underlying abstract.Scheme
	cost       int
}

// An implementation of Scheme implementing Python passlib's `$bcrypt-sha256$`
// bcrypt variant. This is bcrypt with a SHA256 prehash, which removes bcrypt's
// password length limitation.
var Crypter abstract.Scheme

// The recommended cost for bcrypt-sha256. This may change with subsequent releases.
const RecommendedCost = bcrypt.RecommendedCost

func init() {
	Crypter = New(bcrypt.RecommendedCost)
}

// Instantiates a new Scheme implementing bcrypt with the given cost.
//
// The recommended cost is RecommendedCost.
func New(cost int) abstract.Scheme {
	return &scheme{
		underlying: bcrypt.New(cost),
		cost:       cost,
	}
}

func (s *scheme) Hash(password string) (string, error) {
	p := s.prehash(password)
	h, err := s.underlying.Hash(p)
	if err != nil {
		return "", err
	}

	return mangle(h), nil
}

func (s *scheme) Verify(password, hash string) error {
	p := s.prehash(password)
	return s.underlying.Verify(p, demangle(hash))
}

func (s *scheme) prehash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	v := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return v
}

func (s *scheme) SupportsStub(stub string) bool {
	return strings.HasPrefix(stub, "$bcrypt-sha256$") && s.underlying.SupportsStub(demangle(stub))
}

func (s *scheme) NeedsUpdate(stub string) bool {
	return s.underlying.NeedsUpdate(demangle(stub))
}

func (s *scheme) String() string {
	return fmt.Sprintf("bcrypt-sha256(%d)", s.cost)
}

func demangle(stub string) string {
	if strings.HasPrefix(stub, "$bcrypt-sha256$2") {
		parts := strings.Split(stub[15:], "$")
		// 0: 2a,12
		// 1: salt
		// 2: hash
		parts0 := strings.Split(parts[0], ",")
		return "$" + parts0[0] + "$" + fmt.Sprintf("%02s", parts0[1]) + "$" + parts[1] + parts[2]
	} else {
		return stub
	}
}

func mangle(hash string) string {
	parts := strings.Split(hash[1:], "$")
	// 0: 2a
	// 1: rounds
	// 2: salt + hash
	salt := parts[2][0:22]
	h := parts[2][22:]
	return "$bcrypt-sha256$" + parts[0] + "," + parts[1] + "$" + salt + "$" + h
}
