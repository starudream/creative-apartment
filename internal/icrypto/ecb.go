package icrypto

import (
	"github.com/forgoer/openssl"

	"github.com/starudream/creative-apartment/internal/idt"
)

func NewAesEcbPKCS5(key []byte) *xAesEcb {
	return &xAesEcb{key: key, padding: openssl.PKCS5_PADDING}
}

type xAesEcb struct {
	key []byte

	padding string
}

var _ Interface = (*xAesEcb)(nil)

func (x xAesEcb) Encrypt(bytes []byte) (idt.Bytes, error) {
	return openssl.AesECBEncrypt(bytes, x.key, x.padding)
}

func (x xAesEcb) Decrypt(bytes []byte) (idt.Bytes, error) {
	return openssl.AesECBDecrypt(bytes, x.key, x.padding)
}
