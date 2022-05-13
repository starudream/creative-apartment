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

	viper.SetEnvPrefix("sca")
	viper.AutomaticEnv()
	viper.SetDefault("debug", true)

	path := viper.GetString("path")

	log.Info().Msgf("path: %s", path)

	if path != "" {
		viper.SetConfigFile(path)
		_ = viper.ReadInConfig()
		path = filepath.Dir(path)
	} else {
		path = ios.Executable()
	}

	ibolt.Init(filepath.Join(path, config.AppName+".bolt"))
}
