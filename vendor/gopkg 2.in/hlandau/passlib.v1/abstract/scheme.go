package abstract

// The Scheme interface provides an abstract interface to an implementation
// of a particular password hashing scheme. The Scheme generates password
// hashes from passwords, verifies passwords using password hashes, randomly
// generates new stubs and can determines whether it recognises a given
// stub or hash. It may also decide to issue upgrades.
type Scheme interface {
	// Hashes a plaintext UTF-8 password using a modular crypt stub. Returns the
	// hashed password in modular crypt format.
	//
	// A modular crypt stub is a prefix of a hash in modular crypt format which
	// expresses all necessary configuration information, such as salt and
	// iteration count. For example, for sha256-crypt, a valid stub would be:
	//
	//     $5$rounds=6000$salt
	//
	// A full modular crypt hash may also be passed as the stub, in which case
	// the hash is ignored.
	Hash(password string) (string, error)

	// Verifies a plaintext UTF-8 password using a modular crypt hash.  Returns
	// an error if the inputs are malformed or the password does not match.
	Verify(password, hash string) (err error)

	// Returns true iff this crypter supports the given stub.
	SupportsStub(stub string) bool

	// Returns true iff this stub needs an update.
	NeedsUpdate(stub string) bool

	// Make a stub with the configured defaults. The salt is generated randomly.
	//MakeStub() (string, error)
}
