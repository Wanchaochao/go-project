package files

import (
	"crypto/sha1"
	"encoding/base32"
)

var enc = base32.NewEncoding("etif7xdq35nvboszk2h6mu4ycrwajgpl").WithPadding(base32.NoPadding)

func GenFilePath(b []byte) string {
	h := sha1.Sum(b)
	s := enc.EncodeToString(h[:])
	return s[:2] + "/" + s[2:]
}
