package raw

import (
	"encoding/base64"
	"strings"
)

var b64 = base64.RawStdEncoding

func Base64Encode(src []byte) (dst string) {
	dst = b64.EncodeToString(src)
	dst = strings.Replace(dst, "+", ".", -1)
	return
}

func Base64Decode(src string) (dst []byte, err error) {
	src = strings.Replace(src, ".", "+", -1)
	dst, err = b64.DecodeString(src)
	return
}
