package icrypto

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/idt"
)

func TestAESECBP5(t *testing.T) {
	x := NewAesEcbPKCS5([]byte(config.AESECBP5Key))

	src := idt.Bytes("1234567890123456")

	s1, err := x.Encrypt(src)
	require.NoError(t, err)
	require.Equal(t, "ZAygwCVRe79dU5Y3AH9HFt3vcH6wmvXL5Yck2qWlsYI=", s1.Base64String())

	s2, err := x.Decrypt(s1)
	require.NoError(t, err)
	require.Equal(t, src.String(), s2.String())
}
