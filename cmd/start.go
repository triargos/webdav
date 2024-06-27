package cmd

import (
	"github.com/spf13/cobra"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/server"
	"log/slog"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the WebDAV server",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Starting server...")
		slog.Info("Creating user directories...")
		server.CreateUserDirectories()
		slog.Info("Hashing passwords...")
		auth.HashPasswords()
		err := config.Write()
		if err != nil {
			slog.Error("Failed to write configuration to disk:", "error", err.Error())
			os.Exit(1)
		}
		slog.Info("Starting http handler...")
		err = server.StartWebdavServer()
		if err != nil {
			slog.Error("Failed to start http handler:", "error", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
