package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/starudream/creative-apartment/config"
)

var rootCmd = &cobra.Command{
	Use:     config.AppName,
	Short:   config.AppName,
	Version: fmt.Sprintf("%s (%s)", config.VERSION, config.BIDTIME),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
