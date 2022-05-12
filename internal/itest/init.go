package itest

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/ios"
)

func Init() {
	log.Logger = log.Output(ilog.NewConsoleWriter())

	viper.AutomaticEnv()
	path := viper.GetString("SCA_PATH")
	if path != "" {
		path = filepath.Dir(path)
	} else {
		path = ios.Executable()
	}

	ibolt.Init(filepath.Join(path, config.AppName+".bolt"))
}
