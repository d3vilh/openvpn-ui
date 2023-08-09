package abstract

import "crypto/subtle"

// Compares two strings (typicaly password hashes) in a secure, constant-time
// fashion. Returns true iff they are equal.
func SecureCompare(a, b string) bool {
	ab := []byte(a)
	bb := []byte(b)
	return subtle.ConstantTimeCompare(ab, bb) == 1
}
