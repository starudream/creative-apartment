package cmd

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/ios"
)

var path string

func init() {
	cobra.OnInitialize(initLogger, initConfig, initDB)

	rootCmd.PersistentFlags().BoolP("debug", "", false, "(env: SCA_DEBUG) show debug information")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "(env: SCA_PATH) configuration file path")
	viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path"))

	rootCmd.AddCommand(runCmd)
}

func initConfig() {
	viper.SetConfigPermissions(0644)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("sca") // starudream - creative - apartment
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("debug", false)
	viper.SetDefault("log.level", "INFO")

	path = viper.GetString("path")

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

	if path != "" {
		path = filepath.Dir(path)
	}
	if v := viper.ConfigFileUsed(); v != "" {
		path = filepath.Dir(v)
	}
	if path == "" {
		path = filepath.Join(ios.UserHomeDir(), ".config", "starudream")
	}

	viper.SetConfigFile(filepath.Join(path, config.AppName+".yaml"))

	level, err := zerolog.ParseLevel(strings.ToLower(viper.GetString("log.level")))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	debug := viper.GetBool("debug")
	if debug {
		level = zerolog.DebugLevel
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}
	zerolog.SetGlobalLevel(level)

	if level < zerolog.NoLevel {
		log.Logger = log.Output(zerolog.MultiLevelWriter(newConsoleWriter(), newFileWriter()))
		if debug {
			log.Logger = log.Logger.With().Caller().Logger()
		}
	}

	zerolog.DefaultContextLogger = &log.Logger

	log.Info().Msgf("workspace path: %s", path)
}

func initDB() {
	config.SetCustomers(viper.Get("customers"))

	ibolt.Init(filepath.Join(path, config.AppName+".bolt"))

	ierr.CheckErr(ibolt.Update(func(tx *ibolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("config"))
		return err
	}))
}

func newConsoleWriter() io.Writer {
	return &zerolog.ConsoleWriter{
		Out:        colorable.NewColorableStdout(),
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
	}
}

func newFileWriter() io.Writer {
	return &zerolog.ConsoleWriter{
		Out: &lumberjack.Logger{
			Filename:  filepath.Join(path, config.AppName+".log"),
			MaxSize:   100,
			MaxAge:    360,
			LocalTime: true,
		},
		NoColor:    true,
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
	}
}

func initLogger() {
	w := ilog.New(log.Output(newConsoleWriter()), "cfg")
	jww.TRACE = w.WithLevel(zerolog.TraceLevel)
	jww.DEBUG = w.WithLevel(zerolog.DebugLevel)
	jww.INFO = w.WithLevel(zerolog.InfoLevel)
	jww.WARN = w.WithLevel(zerolog.WarnLevel)
	jww.ERROR = w.WithLevel(zerolog.ErrorLevel)
	jww.CRITICAL = w.WithLevel(zerolog.FatalLevel)
	jww.FATAL = w.WithLevel(zerolog.FatalLevel)
	jww.LOG = w.WithLevel(zerolog.TraceLevel)
}
