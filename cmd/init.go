package cmd

import (
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/icfg"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/ios"
	"github.com/starudream/creative-apartment/internal/iseq"
)

func init() {
	cobra.OnInitialize(initLogger, initConfig, initDB)

	rootCmd.PersistentFlags().BoolP("debug", "", false, "(env: SCA_DEBUG) show debug information")
	ierr.CheckErr(viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")))

	rootCmd.PersistentFlags().StringP("path", "", "", "(env: SCA_PATH) configuration file path")
	ierr.CheckErr(viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path")))

	rootCmd.PersistentFlags().IntP("port", "", 8089, "(env: SCA_PORT) http server port")
	ierr.CheckErr(viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port")))

	rootCmd.PersistentFlags().IntP("secret", "", 8089, "(env: SCA_SECRET) login secret")
	ierr.CheckErr(viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret")))
}

func initConfig() {
	viper.SetConfigPermissions(0644)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("sca") // starudream - creative - apartment
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("debug", false)
	viper.SetDefault("log.level", "INFO")

	path := viper.GetString("path")

	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigName(config.AppName)

		viper.AddConfigPath(ios.Executable())

		viper.AddConfigPath(ios.UserHomeDir())
		viper.AddConfigPath(filepath.Join(ios.UserHomeDir(), ".config", "starudream"))
	}

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok && path != "" {
		ierr.CheckErr(err)
	}

	path = viper.ConfigFileUsed()
	if path == "" {
		path = filepath.Join(ios.UserHomeDir(), ".config", "starudream", config.AppName+".yaml")
	}

	viper.SetConfigFile(path)

	level, err := zerolog.ParseLevel(strings.ToLower(viper.GetString("log.level")))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	debug := viper.GetBool("debug")
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	if level < zerolog.NoLevel {
		lfn := filepath.Join(filepath.Dir(viper.GetString("path")), config.AppName+".log")
		log.Logger = log.Output(zerolog.MultiLevelWriter(ilog.NewConsoleWriter(), ilog.NewFileWriter(lfn)))
		if debug {
			log.Logger = log.With().Caller().Logger()
		}
	}

	zerolog.DefaultContextLogger = &log.Logger

	log.Info().Msgf("[cfg] workspace path: %s", filepath.Dir(path))

	secret := viper.GetString("secret")
	if secret == "" {
		secret = iseq.UUID()
		viper.Set("secret", secret)
		icfg.Save()
	}

	log.Info().Msgf("[cfg] login secret: %s", secret)
}

func initDB() {
	config.SetCustomers(viper.Get("customers"))

	ibolt.Init(filepath.Join(filepath.Dir(viper.GetString("path")), config.AppName+".bolt"))

	ierr.CheckErr(ibolt.Update(func(tx *ibolt.Tx) error {
		buckets := []string{"config", "customer"}
		for i := 0; i < len(buckets); i++ {
			_, err := tx.CreateBucketIfNotExists([]byte(buckets[i]))
			if err != nil {
				return err
			}
		}
		return nil
	}))
}

func initLogger() {
	w := ilog.New(log.Output(ilog.NewConsoleWriter()), "cfg")
	jww.TRACE = w.WithLevel(zerolog.TraceLevel)
	jww.DEBUG = w.WithLevel(zerolog.DebugLevel)
	jww.INFO = w.WithLevel(zerolog.InfoLevel)
	jww.WARN = w.WithLevel(zerolog.WarnLevel)
	jww.ERROR = w.WithLevel(zerolog.ErrorLevel)
	jww.CRITICAL = w.WithLevel(zerolog.FatalLevel)
	jww.FATAL = w.WithLevel(zerolog.FatalLevel)
	jww.LOG = w.WithLevel(zerolog.TraceLevel)
}
