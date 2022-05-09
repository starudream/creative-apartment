package idt

import (
	"encoding/base64"
	"encoding/hex"
)

type Bytes []byte

func (bs Bytes) String() string {
	return string(bs)
}

func (bs Bytes) Hex() string {
	return hex.EncodeToString(bs)
}

func (bs Bytes) Base64String() string {
	return base64.StdEncoding.EncodeToString(bs)
}
