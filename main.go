package main

import (
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/cmd"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ierr"
)

func main() {
	defer ierr.Recover(ibolt.Close, ierr.WrapErrFunc(viper.WriteConfig))

	cmd.Execute()
}
