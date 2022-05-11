package main

import (
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/cmd"
	"github.com/starudream/creative-apartment/internal/app"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/igin"
)

func main() {
	defer ierr.Recover(app.Stop, igin.Close, ibolt.Close, ierr.WrapErrFunc(viper.WriteConfig))

	cmd.Execute()

	ierr.CheckErr(app.Go())
}
