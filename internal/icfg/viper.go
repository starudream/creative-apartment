package icfg

import (
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/internal/ierr"
)

var (
	keys = []string{"secret", "customers", "dingtalk"}

	done int32
)

func Done() {
	atomic.StoreInt32(&done, 1)
}

func Save() {
	if atomic.LoadInt32(&done) != 1 {
		return
	}

	ierr.CheckErr(os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), 0755))
	nViper := viper.New()
	for i := 0; i < len(keys); i++ {
		nViper.Set(keys[i], viper.Get(keys[i]))
	}
	ierr.CheckErr(nViper.WriteConfigAs(viper.ConfigFileUsed()))
}
