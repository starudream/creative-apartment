package main

import (
	"github.com/starudream/creative-apartment/cmd"
	"github.com/starudream/creative-apartment/internal/app"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/icfg"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/igin"
)

func main() {
	app.Recover(icfg.Save, ibolt.Close, igin.Close)

	cmd.Execute()

	ierr.CheckErr(app.Go())
}
