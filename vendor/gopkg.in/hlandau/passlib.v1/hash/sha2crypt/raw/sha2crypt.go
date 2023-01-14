// Package raw provides a raw implementation of the sha256-crypt and sha512-crypt primitives.
package raw

import "io"
import "fmt"
import "hash"
import "crypto/sha256"
import "crypto/sha512"

// The minimum number of rounds permissible for sha256-crypt and sha512-crypt.
const MinimumRounds = 1000

// The maximum number of rounds permissible for sha256-crypt and sha512-crypt.
// Don't use this!
const MaximumRounds = 999999999

// This is the 'default' number of rounds for sha256-crypt and sha512-crypt. If
// this rounds value is used the number of rounds is not explicitly specified
// in the modular crypt format, as it is the default.
const DefaultRounds = 5000

// This is the recommended number of rounds for sha256-crypt and sha512-crypt.
// This may change with subsequent releases of this package. It is recommended
// that you invoke sha256-crypt or sha512-crypt with this value, or a value
// proportional to it.
const RecommendedRounds = 10000

// Calculates sha256-crypt. The password must be in plaintext and be a UTF-8
// string.
//
// The salt must be a valid ASCII between 0 and 16 characters in length
// inclusive.
//
// See the constants in this package for suggested values for rounds.
//
// Rounds must be in the range 1000 <= rounds <= 999999999. The function panics
// if this is not the case.
//
// The output is in modular crypt format.
func Crypt256(password, salt string, rounds int) string {
	return "$5" + shaCrypt(password, salt, rounds, sha256.New, transpose256)
}

// Calculates sha256-crypt. The password must be in plaintext and be a UTF-8
// string.
//
// The salt must be a valid ASCII between 0 and 16 characters in length
// inclusive.
//
// See the constants in this package for suggested values for rounds.
//
// Rounds must be in the range 1000 <= rounds <= 999999999. The function panics
// if this is not the case.
//
// The output is in modular crypt format.
func Crypt512(password, salt string, rounds int) string {
	return "$6" + shaCrypt(password, salt, rounds, sha512.New, transpose512)
}

func shaCrypt(password, salt string, rounds int, newHash func() hash.Hash, transpose func(b []byte)) string {
	if rounds < MinimumRounds || rounds > MaximumRounds {
		panic("sha256-crypt rounds must be in 1000 <= rounds <= 999999999")
	}

	passwordb := []byte(password)
	saltb := []byte(salt)
	if len(saltb) > 16 {
		panic("salt must not exceed 16 bytes")
	}

	// B
	b := newHash()
	b.Write(passwordb)
	b.Write(saltb)
	b.Write(passwordb)
	bsum := b.Sum(nil)

	// A
	a := newHash()
	a.Write(passwordb)
	a.Write(saltb)
	repeat(a, bsum, len(passwordb))

	plen := len(passwordb)
	for plen != 0 {
		if (plen & 1) != 0 {
			a.Write(bsum)
		} else {
			a.Write(passwordb)
		}
		plen = plen >> 1
	}

	asum := a.Sum(nil)

	// DP
	dp := newHash()
	for i := 0; i < len(passwordb); i++ {
		dp.Write(passwordb)
	}

	dpsum := dp.Sum(nil)

	// P
	p := make([]byte, len(passwordb))
	repeatTo(p, dpsum)

	// DS
	ds := newHash()
	for i := 0; i < (16 + int(asum[0])); i++ {
		ds.Write(saltb)
	}

	dssum := ds.Sum(nil)[0:len(saltb)]

	// S
	s := make([]byte, len(saltb))
	repeatTo(s, dssum)

	// C
	cur := asum[:]
	for i := 0; i < rounds; i++ {
		c := newHash()
		if (i & 1) != 0 {
			c.Write(p)
		} else {
			c.Write(cur)
		}
		if (i % 3) != 0 {
			c.Write(s)
		}
		if (i % 7) != 0 {
			c.Write(p)
		}
		if (i & 1) == 0 {
			c.Write(p)
		} else {
			c.Write(cur)
		}
		cur = c.Sum(nil)[:]
	}

	// Transposition
	transpose(cur)

	// Hash
	hstr := EncodeBase64(cur)

	if rounds == DefaultRounds {
		return fmt.Sprintf("$%s$%s", salt, hstr)
	}

	return fmt.Sprintf("$rounds=%d$%s$%s", rounds, salt, hstr)
}

func repeat(w io.Writer, b []byte, sz int) {
	var i int
	for i = 0; (i + len(b)) <= sz; i += len(b) {
		w.Write(b)
	}
	w.Write(b[0 : sz-i])
}

func repeatTo(out []byte, b []byte) {
	if len(b) == 0 {
		return
	}

	var i int
	for i = 0; (i + len(b)) <= len(out); i += len(b) {
		copy(out[i:], b)
	}
	copy(out[i:], b)
}

func transpose256(b []byte) {
	b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19], b[20], b[21], b[22], b[23], b[24], b[25], b[26], b[27], b[28], b[29], b[30], b[31] =
		b[20], b[10], b[0], b[11], b[1], b[21], b[2], b[22], b[12], b[23], b[13], b[3], b[14], b[4], b[24], b[5], b[25], b[15], b[26], b[16], b[6], b[17], b[7], b[27], b[8], b[28], b[18], b[29], b[19], b[9], b[30], b[31]
}

func transpose512(b []byte) {
	b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19], b[20], b[21], b[22], b[23], b[24], b[25], b[26], b[27], b[28], b[29], b[30], b[31], b[32], b[33], b[34], b[35], b[36], b[37], b[38], b[39], b[40], b[41], b[42], b[43], b[44], b[45], b[46], b[47], b[48], b[49], b[50], b[51], b[52], b[53], b[54], b[55], b[56], b[57], b[58], b[59], b[60], b[61], b[62], b[63] =
		b[42], b[21], b[0], b[1], b[43], b[22], b[23], b[2], b[44], b[45], b[24], b[3], b[4], b[46], b[25], b[26], b[5], b[47], b[48], b[27], b[6], b[7], b[49], b[28], b[29], b[8], b[50], b[51], b[30], b[9], b[10], b[52], b[31], b[32], b[11], b[53], b[54], b[33], b[12], b[13], b[55], b[34], b[35], b[14], b[56], b[57], b[36], b[15], b[16], b[58], b[37], b[38], b[17], b[59], b[60], b[39], b[18], b[19], b[61], b[40], b[41], b[20], b[62], b[63]
}

// © 2008-2012 Assurance Technologies LLC.  (Python passlib)  BSD License
// © 2014 Hugo Landau <hlandau@devever.net>  BSD License
