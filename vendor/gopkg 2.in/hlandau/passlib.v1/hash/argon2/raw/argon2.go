// Package raw provides a raw implementation of the modular-crypt-wrapped Argon2i primitive.
package raw

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strconv"
	"strings"
)

// The current recommended time value for interactive logins.
const RecommendedTime uint32 = 4

// The current recommended memory for interactive logins.
const RecommendedMemory uint32 = 32 * 1024

// The current recommended number of threads for interactive logins.
const RecommendedThreads uint8 = 4

// Wrapper for golang.org/x/crypto/argon2 implementing a sensible
// hashing interface.
//
// password should be a UTF-8 plaintext password.
// salt should be a random salt value in binary form.
//
// Time, memory, and threads are parameters to argon2.
//
// Returns an argon2 encoded hash.
func Argon2(password string, salt []byte, time, memory uint32, threads uint8) string {
	passwordb := []byte(password)

	hash := argon2.Key(passwordb, salt, time, memory, threads, 32)

	hstr := base64.RawStdEncoding.EncodeToString(hash)
	sstr := base64.RawStdEncoding.EncodeToString(salt)

	return fmt.Sprintf("$argon2i$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, sstr, hstr)
}

// Indicates that a password hash or stub is invalid.
var ErrInvalidStub = fmt.Errorf("invalid argon2 password stub")

// Indicates that a key-value pair in the configuration part is malformed.
var ErrInvalidKeyValuePair = fmt.Errorf("invalid argon2 key-value pair")

// Indicates that the version part had the wrong number of parameters.
var ErrParseVersion = fmt.Errorf("version section has wrong number of parameters")

// Indicates that the hash config part had the wrong number of parameters.
var ErrParseConfig = fmt.Errorf("hash config section has wrong number of parameters")

// Indicates that the version parameter ("v") was missing in the version part,
// even though it is required.
var ErrMissingVersion = fmt.Errorf("version parameter (v) is missing")

// Indicates that the memory parameter ("m") was mossing in the hash config
// part, even though it is required.
var ErrMissingMemory = fmt.Errorf("memory parameter (m) is missing")

// Indicates that the time parameter ("t") was mossing in the hash config part,
// even though it is required.
var ErrMissingTime = fmt.Errorf("time parameter (t) is missing")

// Indicates that the parallelism parameter ("p") was mossing in the hash config
// part, even though it is required.
var ErrMissingParallelism = fmt.Errorf("parallelism parameter (p) is missing")

// Parses an argon2 encoded hash.
//
// The format is as follows:
//
//   $argon2i$v=version$m=memory,t=time,p=threads$salt$hash   // hash
//   $argon2i$v=version$m=memory,t=time,p=threads$salt        // stub
//
func Parse(stub string) (salt, hash []byte, version int, time, memory uint32, parallelism uint8, err error) {
	if len(stub) < 26 || !strings.HasPrefix(stub, "$argon2i$") {
		err = ErrInvalidStub
		return
	}

	// $argon2i$  v=version$m=memory,t=time,p=threads$salt-base64$hash-base64
	parts := strings.Split(stub[9:], "$")

	// version-params$hash-config-params$salt[$hash]
	if len(parts) < 3 || len(parts) > 4 {
		err = ErrInvalidStub
		return
	}

	// Parse the first configuration part, the version parameters.
	versionParams, err := parseKeyValuePair(parts[0])
	if err != nil {
		return
	}

	// Must be exactly one parameter in the version part.
	if len(versionParams) != 1 {
		err = ErrParseVersion
		return
	}

	// It must be "v".
	val, ok := versionParams["v"]
	if !ok {
		err = ErrMissingVersion
		return
	}

	version = int(val)

	// Parse the second configuration part, the hash config parameters.
	hashParams, err := parseKeyValuePair(parts[1])
	if err != nil {
		return
	}

	// It must have exactly three parameters.
	if len(hashParams) != 3 {
		err = ErrParseConfig
		return
	}

	// Memory parameter.
	val, ok = hashParams["m"]
	if !ok {
		err = ErrMissingMemory
		return
	}

	memory = uint32(val)

	// Time parameter.
	val, ok = hashParams["t"]
	if !ok {
		err = ErrMissingTime
		return
	}

	time = uint32(val)

	// Parallelism parameter.
	val, ok = hashParams["p"]
	if !ok {
		err = ErrMissingParallelism
		return
	}

	parallelism = uint8(val)

	// Decode salt.
	salt, err = base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return
	}

	// Decode hash if present.
	if len(parts) >= 4 {
		hash, err = base64.RawStdEncoding.DecodeString(parts[3])
	}

	return
}

func parseKeyValuePair(pairs string) (result map[string]uint64, err error) {
	result = map[string]uint64{}

	parameterParts := strings.Split(pairs, ",")

	for _, parameter := range parameterParts {
		parts := strings.SplitN(parameter, "=", 2)
		if len(parts) != 2 {
			err = ErrInvalidKeyValuePair
			return
		}

		parsedi, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			return result, err
		}

		result[parts[0]] = parsedi
	}

	return result, nil
}
