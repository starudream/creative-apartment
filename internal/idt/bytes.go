package idt

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/spf13/cast"
)

type Bytes []byte

func (bs Bytes) String() string {
	return string(bs)
}

func (bs Bytes) Hex() string {
	return hex.EncodeToString(bs)
}

func (bs Bytes) Base64Std() string {
	return base64.StdEncoding.EncodeToString(bs)
}

func (bs Bytes) Base64URL() string {
	return base64.URLEncoding.EncodeToString(bs)
}

func ToBytes(v any) Bytes {
	return Bytes(cast.ToString(v))
}
