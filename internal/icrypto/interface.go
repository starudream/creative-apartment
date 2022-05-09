package icrypto

import (
	"github.com/starudream/creative-apartment/internal/idt"
)

type Interface interface {
	Encrypt(bytes []byte) (idt.Bytes, error)
	Decrypt(bytes []byte) (idt.Bytes, error)
}
