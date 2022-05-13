package ibot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDingtalk_SendMessage(t *testing.T) {
	require.NoError(t, Dingtalk.SendMessage("【2022-05-13】\n耗电 30.00，电费 30.00，剩余 30.00\n耗水 30.00，水费 30.00，剩余 30.00"))
}
