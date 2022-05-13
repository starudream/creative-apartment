package ibot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDingtalk_SendMessage(t *testing.T) {
	require.NoError(t, Dingtalk.SendMessage(""))
}
