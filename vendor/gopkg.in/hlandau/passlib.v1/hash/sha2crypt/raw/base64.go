package raw

const bmap = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Encodes a byte string using the sha2-crypt base64 variant.
func EncodeBase64(b []byte) string {
	o := make([]byte, len(b)/3*4+4)

	for i, j := 0, 0; i < len(b); {
		b1 := b[i]
		b2 := byte(0)
		b3 := byte(0)

		if (i + 1) < len(b) {
			b2 = b[i+1]
		}
		if (i + 2) < len(b) {
			b3 = b[i+2]
		}

		o[j] = bmap[(b1 & 0x3F)]
		o[j+1] = bmap[((b1&0xC0)>>6)|((b2&0x0F)<<2)]
		o[j+2] = bmap[((b2&0xF0)>>4)|((b3&0x03)<<4)]
		o[j+3] = bmap[(b3&0xFC)>>2]
		i += 3
		j += 4
	}

	s := string(o)
	return s[0 : len(b)*4/3-(len(b)%4)+1]
}

// © 2008-2012 Assurance Technologies LLC.  (Python passlib)  BSD License
// © 2014 Hugo Landau <hlandau@devever.net>  BSD License
