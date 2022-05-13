package ibot

import (
	"os"
	"testing"

	"github.com/starudream/creative-apartment/internal/itest"
)

func TestMain(m *testing.M) {
	itest.Init()

	os.Exit(m.Run())
}
