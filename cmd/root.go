package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webdav",
	Short: "The webdav server for your files",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if !fs.PathExists("/etc/webdav") {
		err := os.Mkdir("/etc/webdav", 0755)
		if err != nil {
			slog.Error("Failed to create config directory:", "error", err.Error())
			os.Exit(1)
		}
	}
	err := config.Read()
	if err != nil {
		slog.Error("Failed to read config file:", "error", err.Error())
		os.Exit(1)
	}
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
