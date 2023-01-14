// Package scrypt implements the scrypt password hashing mechanism, wrapped in
// the modular crypt format.
package scrypt

import "fmt"
import "expvar"
import "strings"
import "crypto/rand"
import "encoding/base64"
import "gopkg.in/hlandau/passlib.v1/hash/scrypt/raw"
import "gopkg.in/hlandau/passlib.v1/abstract"

var cScryptSHA256HashCalls = expvar.NewInt("passlib.scryptsha256.hashCalls")
var cScryptSHA256VerifyCalls = expvar.NewInt("passlib.scryptsha256.verifyCalls")

// An implementation of Scheme performing scrypt-sha256.
//
// Uses the recommended values for N,r,p defined in raw.
var SHA256Crypter abstract.Scheme

func init() {
	SHA256Crypter = NewSHA256(
		raw.RecommendedN,
		raw.Recommendedr,
		raw.Recommendedp,
	)
}

// Returns an implementation of Scheme implementing scrypt-sha256
// with the specified parameters.
func NewSHA256(N, r, p int) abstract.Scheme {
	return &scryptSHA256Crypter{
		nN: N,
		r:  r,
		p:  p,
	}
}

type scryptSHA256Crypter struct {
	nN, r, p int
}

func (c *scryptSHA256Crypter) SetParams(N, r, p int) error {
	c.nN = N
	c.r = r
	c.p = p
	return nil
}

func (c *scryptSHA256Crypter) SupportsStub(stub string) bool {
	return strings.HasPrefix(stub, "$s2$")
}

func (c *scryptSHA256Crypter) Hash(password string) (string, error) {
	cScryptSHA256HashCalls.Add(1)

	stub, err := c.makeStub()
	if err != nil {
		return "", err
	}

	_, newHash, _, _, _, _, err := c.hash(password, stub)
	return newHash, err
}

func (c *scryptSHA256Crypter) Verify(password, hash string) (err error) {
	cScryptSHA256VerifyCalls.Add(1)

	_, newHash, _, _, _, _, err := c.hash(password, hash)
	if err == nil && !abstract.SecureCompare(hash, newHash) {
		err = abstract.ErrInvalidPassword
	}

	return
}

func (c *scryptSHA256Crypter) NeedsUpdate(stub string) bool {
	salt, _, N, r, p, err := raw.Parse(stub)
	if err != nil {
		return false // ...
	}

	return c.needsUpdate(salt, N, r, p)
}

func (c *scryptSHA256Crypter) needsUpdate(salt []byte, N, r, p int) bool {
	return len(salt) < 18 || N < c.nN || r < c.r || p < c.p
}

func (c *scryptSHA256Crypter) hash(password, stub string) (oldHashRaw []byte, newHash string, salt []byte, N, r, p int, err error) {
	salt, oldHashRaw, N, r, p, err = raw.Parse(stub)
	if err != nil {
		return
	}

	return oldHashRaw, raw.ScryptSHA256(password, salt, N, r, p), salt, N, r, p, nil
}

func (c *scryptSHA256Crypter) makeStub() (string, error) {
	buf := make([]byte, 18)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	salt := base64.StdEncoding.EncodeToString(buf)

	return fmt.Sprintf("$s2$%d$%d$%d$%s", c.nN, c.r, c.p, salt), nil
}

func (c *scryptSHA256Crypter) String() string {
	return fmt.Sprintf("scrypt-sha256(%d,%d,%d)", c.nN, c.r, c.p)
}
