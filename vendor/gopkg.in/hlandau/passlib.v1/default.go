package passlib

import (
	"fmt"
	"gopkg.in/hlandau/passlib.v1/abstract"
	"gopkg.in/hlandau/passlib.v1/hash/argon2"
	"gopkg.in/hlandau/passlib.v1/hash/bcrypt"
	"gopkg.in/hlandau/passlib.v1/hash/bcryptsha256"
	"gopkg.in/hlandau/passlib.v1/hash/pbkdf2"
	"gopkg.in/hlandau/passlib.v1/hash/scrypt"
	"gopkg.in/hlandau/passlib.v1/hash/sha2crypt"
	"time"
)

// This is the first and default set of defaults used by passlib. It prefers
// scrypt-sha256. It is now obsolete.
const Defaults20160922 = "20160922"

// This is the most up-to-date set of defaults preferred by passlib. It prefers
// Argon2i. You must opt into it by calling UseDefaults at startup.
const Defaults20180601 = "20180601"

// This value, when passed to UseDefaults, causes passlib to always use the
// very latest set of defaults. DO NOT use this unless you are sure that
// opportunistic hash upgrades will not cause breakage for your application
// when future versions of passlib are released. See func UseDefaults.
const DefaultsLatest = "latest"

// Default schemes as of 2016-09-22.
var defaultSchemes20160922 = []abstract.Scheme{
	scrypt.SHA256Crypter,
	argon2.Crypter,
	sha2crypt.Crypter512,
	sha2crypt.Crypter256,
	bcryptsha256.Crypter,
	pbkdf2.SHA512Crypter,
	pbkdf2.SHA256Crypter,
	bcrypt.Crypter,
	pbkdf2.SHA1Crypter,
}

// Default schemes as of 2018-06-01.
var defaultSchemes20180601 = []abstract.Scheme{
	argon2.Crypter,
	scrypt.SHA256Crypter,
	sha2crypt.Crypter512,
	sha2crypt.Crypter256,
	bcryptsha256.Crypter,
	pbkdf2.SHA512Crypter,
	pbkdf2.SHA256Crypter,
	bcrypt.Crypter,
	pbkdf2.SHA1Crypter,
}

// The default schemes, most preferred first. The first scheme will be used to
// hash passwords, and any of the schemes may be used to verify existing
// passwords. The contents of this value may change with subsequent releases.
//
// If you want to change this, set DefaultSchemes to a slice to an
// abstract.Scheme array of your own construction, rather than mutating the
// array the slice points to.
//
// To see the default schemes used in the current release of passlib, see
// default.go. See also the UseDefaults function for more information on how
// the list of default schemes is determined. The default value of
// DefaultSchemes (the default defaults) won't change; you need to call
// UseDefaults to allow your application to upgrade to newer hashing schemes
// (or set DefaultSchemes manually, or create a custom context with its own
// schemes set).
var DefaultSchemes []abstract.Scheme

func init() {
	DefaultSchemes = defaultSchemes20160922
}

// It is strongly recommended that you call this function like this before using passlib:
//
//   passlib.UseDefaults("YYYYMMDD")
//
// where YYYYMMDD is a date. This will be used to select the preferred scheme
// to use. If you do not call UseDefaults, the preferred scheme (the first item
// in the default schemes list) current as of 2016-09-22 will always be used,
// meaning that upgrade will not occur even though better schemes are now
// available.
//
// Note that even if you don't call this function, new schemes will still be
// added to DefaultSchemes over time as non-initial values (items not at index
// 0), so servers will always, by default, be able to validate all schemes
// which passlib supports at any given time.
//
// The reason you must call this function is as follows: If passlib is deployed
// as part of a web application in a multi-server deployment, and passlib is
// updated, and the new version of that application with the updated passlib is
// deployed, that upgrade process is unlikely to be instantaneous. Old versions
// of the web application may continue to run on some servers. If merely
// upgrading passlib caused password hashes to be upgraded to the newer scheme
// on login, the older daemons may not be able to validate these passwords and
// users may have issues logging in. Although this can be ameliorated to some
// extent by introducing a new scheme to passlib, waiting some months, and only
// then making this the default, this could still cause issued if passlib is
// only updated very occasionally.
//
// Thus, you should update your call to UseDefaults only when all servers have
// been upgraded, and it is thus guaranteed that they will all be able to
// verify the new scheme. Making this value loadable from a configuration file
// is recommended.
//
// If you are using a single-server configuration, you can use the special
// value "latest" here (or, equivalently, a date far into the future), which
// will always use the most preferred scheme. This is hazardous in a
// multi-server environment.
//
// The constants beginning 'Defaults' in this package document dates
// which are meaningful to this function. The constant values they are equal to
// will never change, so there is no need to use them instead of string
// literals, although you may if you wish; they are intended mainly as
// documentation as to the significance of various dates.
//
// Example for opting in to the latest set of defaults:
//
//   passlib.UseDefaults(passlib.Defaults20180601)
//
func UseDefaults(date string) error {
	if date == "latest" {
		DefaultSchemes = defaultSchemes20180601
		return nil
	}

	t, err := time.ParseInLocation("20060102", date, time.UTC)
	if err != nil {
		return fmt.Errorf("invalid time string passed to passlib.UseDefaults: %q", date)
	}

	if !t.Before(time.Date(2016, 9, 22, 0, 0, 0, 0, time.UTC)) {
		DefaultSchemes = defaultSchemes20180601
		return nil
	}

	DefaultSchemes = defaultSchemes20160922
	return nil
}
