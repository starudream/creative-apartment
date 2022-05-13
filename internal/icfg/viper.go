package icfg

import (
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/internal/ierr"
)

var keys = []string{"secret", "customers"}

func Save() {
	nViper := viper.New()
	for i := 0; i < len(keys); i++ {
		nViper.Set(keys[i], viper.Get(keys[i]))
	}
	ierr.CheckErr(nViper.WriteConfigAs(viper.ConfigFileUsed()))
}
