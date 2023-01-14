// Package argon2 implements the argon2 password hashing mechanism, wrapped in
// the argon2 encoded format.
package argon2

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"gopkg.in/hlandau/passlib.v1/abstract"
	"gopkg.in/hlandau/passlib.v1/hash/argon2/raw"
)

// An implementation of Scheme performing argon2 hashing.
//
// Uses the recommended values for time, memory and threads defined in raw.
var Crypter abstract.Scheme

const saltLength = 16

func init() {
	Crypter = New(
		raw.RecommendedTime,
		raw.RecommendedMemory,
		raw.RecommendedThreads,
	)
}

// Returns an implementation of Scheme implementing argon2
// with the specified parameters.
func New(time, memory uint32, threads uint8) abstract.Scheme {
	return &scheme{
		time:    time,
		memory:  memory,
		threads: threads,
	}
}

type scheme struct {
	time, memory uint32
	threads      uint8
}

func (c *scheme) SetParams(time, memory uint32, threads uint8) error {
	c.time = time
	c.memory = memory
	c.threads = threads
	return nil
}

func (c *scheme) SupportsStub(stub string) bool {
	return strings.HasPrefix(stub, "$argon2i$")
}

func (c *scheme) Hash(password string) (string, error) {

	stub, err := c.makeStub()
	if err != nil {
		return "", err
	}

	_, newHash, _, _, _, _, _, err := c.hash(password, stub)
	return newHash, err
}

func (c *scheme) Verify(password, hash string) (err error) {

	_, newHash, _, _, _, _, _, err := c.hash(password, hash)
	if err == nil && !abstract.SecureCompare(hash, newHash) {
		err = abstract.ErrInvalidPassword
	}

	return
}

func (c *scheme) NeedsUpdate(stub string) bool {
	salt, _, version, time, memory, threads, err := raw.Parse(stub)
	if err != nil {
		return false // ...
	}

	return c.needsUpdate(salt, version, time, memory, threads)
}

func (c *scheme) needsUpdate(salt []byte, version int, time, memory uint32, threads uint8) bool {
	return len(salt) < saltLength || version < argon2.Version || time < c.time || memory < c.memory || threads < c.threads
}

func (c *scheme) hash(password, stub string) (oldHashRaw []byte, newHash string, salt []byte, version int, memory, time uint32, threads uint8, err error) {

	salt, oldHashRaw, version, time, memory, threads, err = raw.Parse(stub)
	if err != nil {
		return
	}

	return oldHashRaw, raw.Argon2(password, salt, time, memory, threads), salt, version, memory, time, threads, nil
}

func (c *scheme) makeStub() (string, error) {
	buf := make([]byte, saltLength)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	salt := base64.RawStdEncoding.EncodeToString(buf)

	return fmt.Sprintf("$argon2i$v=%d$m=%d,t=%d,p=%d$%s$", argon2.Version, c.memory, c.time, c.threads, salt), nil
}

func (c *scheme) String() string {
	return fmt.Sprintf("argon2(%d,%d,%d,%d)", argon2.Version, c.memory, c.time, c.threads)
}
