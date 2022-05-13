package icrypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXECB_SHA256(t *testing.T) {
	require.Equal(t, "AGGhY3zln8Ha0aUnnnw86Qh8jLd2E+fwMDvbDEaYl5A=", ECB.SHA256("1234567890123456", "1234567890").Base64Std())
}
