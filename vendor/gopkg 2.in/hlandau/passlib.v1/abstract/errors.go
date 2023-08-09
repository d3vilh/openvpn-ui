// Package abstract contains the abstract description of the Scheme interface,
// plus supporting error definitions.
package abstract

import "fmt"

// Indicates that password verification failed because the provided password
// does not match the provided hash.
var ErrInvalidPassword = fmt.Errorf("invalid password")

// Indicates that password verification is not possible because the hashing
// scheme used by the hash provided is not supported.
var ErrUnsupportedScheme = fmt.Errorf("unsupported scheme")

// © 2014 Hugo Landau <hlandau@devever.net>  MIT License
