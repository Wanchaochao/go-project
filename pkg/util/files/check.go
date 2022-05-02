package files

import (
	"net/http"
)

func CheckImage(b []byte) (ext string, ok bool) {
	mine := http.DetectContentType(b)
	if mine == "image/jpeg" || mine == "image/png" || mine == "image/gif" {
		ext = mine[6:]
		ok = true
	}
	return
}
