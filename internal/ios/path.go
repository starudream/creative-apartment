package ios

import (
	"os"

	"github.com/starudream/creative-apartment/internal/ierr"
)

var (
	exec string
	home string
)

func init() {
	var err error

	exec, err = os.Executable()
	ierr.CheckErr(err)

	home, err = os.UserHomeDir()
	ierr.CheckErr(err)
}

func Executable() string {
	return exec
}

func UserHomeDir() string {
	return home
}
