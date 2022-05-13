package icrypto

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/starudream/creative-apartment/internal/idt"
)

var ECB = xECB{}

type xECB struct {
}

func (x xECB) SHA256(src, key string) idt.Bytes {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	return h.Sum(nil)
}
