package itest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/ios"
)

func Init(m *testing.M) {
	log.Logger = log.Output(ilog.NewConsoleWriter())

	viper.SetEnvPrefix("sca")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
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

	os.Exit(m.Run())
}
