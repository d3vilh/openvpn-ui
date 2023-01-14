// Package passlib provides a simple password hashing and verification
// interface abstracting multiple password hashing schemes.
//
// After initialisation, most people need concern themselves only with the
// functions Hash and Verify, which uses the default context and sensible
// defaults.
//
// Library Initialization
//
// You should initialise the library before using it with the following line.
//
//   // Call this at application startup.
//   passlib.UseDefaults(passlib.Defaults20180601)
//
// See func UseDefaults for details.
package passlib // import "gopkg.in/hlandau/passlib.v1"

import (
	"gopkg.in/hlandau/easymetric.v1/cexp"
	"gopkg.in/hlandau/passlib.v1/abstract"
)

var cHashCalls = cexp.NewCounter("passlib.ctx.hashCalls")
var cVerifyCalls = cexp.NewCounter("passlib.ctx.verifyCalls")
var cSuccessfulVerifyCalls = cexp.NewCounter("passlib.ctx.successfulVerifyCalls")
var cFailedVerifyCalls = cexp.NewCounter("passlib.ctx.failedVerifyCalls")
var cSuccessfulVerifyCallsWithUpgrade = cexp.NewCounter("passlib.ctx.successfulVerifyCallsWithUpgrade")
var cSuccessfulVerifyCallsDeferringUpgrade = cexp.NewCounter("passlib.ctx.successfulVerifyCallsDeferringUpgrade")

// A password hashing context, that uses a given set of schemes to hash and
// verify passwords.
type Context struct {
	// Slice of schemes to use, most preferred first.
	//
	// If left uninitialized, a sensible default set of schemes will be used.
	//
	// An upgrade hash (see the newHash return value of the Verify method of the
	// abstract.Scheme interface) will be issued whenever a password is validated
	// using a scheme which is not the first scheme in this slice.
	Schemes []abstract.Scheme
}

func (ctx *Context) schemes() []abstract.Scheme {
	if ctx.Schemes == nil {
		return DefaultSchemes
	}

	return ctx.Schemes
}

// Hashes a UTF-8 plaintext password using the context and produces a password hash.
//
// If stub is "", one is generated automaticaly for the preferred password hashing
// scheme; you should specify stub as "" in almost all cases.
//
// The provided or randomly generated stub is used to deterministically hash
// the password. The returned hash is in modular crypt format.
//
// If the context has not been specifically configured, a sensible default policy
// is used. See the fields of Context.
func (ctx *Context) Hash(password string) (hash string, err error) {
	cHashCalls.Add(1)

	return ctx.schemes()[0].Hash(password)
}

// Verifies a UTF-8 plaintext password using a previously derived password hash
// and the default context. Returns nil err only if the password is valid.
//
// If the hash is determined to be deprecated based on the context policy, and
// the password is valid, the password is hashed using the preferred password
// hashing scheme and returned in newHash. You should use this to upgrade any
// stored password hash in your database.
//
// newHash is empty if the password was not valid or if no upgrade is required.
//
// You should treat any non-nil err as a password verification error.
func (ctx *Context) Verify(password, hash string) (newHash string, err error) {
	return ctx.verify(password, hash, true)
}

// Like Verify, but does not hash an upgrade password when upgrade is required.
func (ctx *Context) VerifyNoUpgrade(password, hash string) error {
	_, err := ctx.verify(password, hash, false)
	return err
}

func (ctx *Context) verify(password, hash string, canUpgrade bool) (newHash string, err error) {
	cVerifyCalls.Add(1)

	for i, scheme := range ctx.schemes() {
		if !scheme.SupportsStub(hash) {
			continue
		}

		err = scheme.Verify(password, hash)
		if err != nil {
			cFailedVerifyCalls.Add(1)
			return "", err
		}

		cSuccessfulVerifyCalls.Add(1)
		if i != 0 || scheme.NeedsUpdate(hash) {
			if canUpgrade {
				cSuccessfulVerifyCallsWithUpgrade.Add(1)

				// If the scheme is not the first scheme, try and rehash with the
				// preferred scheme.
				if newHash, err2 := ctx.Hash(password); err2 == nil {
					return newHash, nil
				}
			} else {
				cSuccessfulVerifyCallsDeferringUpgrade.Add(1)
			}
		}

		return "", nil
	}

	return "", abstract.ErrUnsupportedScheme
}

// Determines whether a stub or hash needs updating according to the policy of
// the context.
func (ctx *Context) NeedsUpdate(stub string) bool {
	for i, scheme := range ctx.schemes() {
		if scheme.SupportsStub(stub) {
			return i != 0 || scheme.NeedsUpdate(stub)
		}
	}

	return false
}

// The default context, which uses sensible defaults. Most users should not
// reconfigure this. The defaults may change over time, so you may wish
// to reconfigure the context or use a custom context if you want precise
// control over the hashes used.
var DefaultContext Context

// Hashes a UTF-8 plaintext password using the default context and produces a
// password hash. Chooses the preferred password hashing scheme based on the
// configured policy. The default policy is sensible.
func Hash(password string) (hash string, err error) {
	return DefaultContext.Hash(password)
}

// Verifies a UTF-8 plaintext password using a previously derived password hash
// and the default context. Returns nil err only if the password is valid.
//
// If the hash is determined to be deprecated based on policy, and the password
// is valid, the password is hashed using the preferred password hashing scheme
// and returned in newHash. You should use this to upgrade any stored password
// hash in your database.
//
// newHash is empty if the password was invalid or no upgrade is required.
//
// You should treat any non-nil err as a password verification error.
func Verify(password, hash string) (newHash string, err error) {
	return DefaultContext.Verify(password, hash)
}

// Like Verify, but never upgrades.
func VerifyNoUpgrade(password, hash string) error {
	return DefaultContext.VerifyNoUpgrade(password, hash)
}

// Uses the default context to determine whether a stub or hash needs updating.
func NeedsUpdate(stub string) bool {
	return DefaultContext.NeedsUpdate(stub)
}

// © 2008-2012 Assurance Technologies LLC.  (Python passlib)  BSD License
// © 2014 Hugo Landau <hlandau@devever.net>  BSD License
