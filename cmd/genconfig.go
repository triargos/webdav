package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var genconfigCmd = &cobra.Command{
	Use:   "genconfig",
	Short: "Generates a new default config file at the config location.",
	Long:  "Generates a new default config file at the config location. If a config file exists, it will only be overridden when the --reset flag is set to true",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Generating config file")
		reset, _ := cmd.Flags().GetBool("reset")
		if reset {
			err := config.WriteDefaultConfig()
			if err != nil {
				slog.Error("Failed to write default config file:", "error", err.Error())
				os.Exit(1)
			}
		} else {
			err := config.Read()
			if err != nil {
				slog.Error("Failed to read config file:", "error", err.Error())
				os.Exit(1)
			}
		}
		slog.Info("Config file generated successfully")
	},
}

func init() {
	rootCmd.AddCommand(genconfigCmd)
	genconfigCmd.Flags().Bool("reset", false, "Reset the config file to the default values")

}
