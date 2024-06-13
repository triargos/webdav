package cmd

import (
	"github.com/spf13/cobra"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/logging"
	"github.com/triargos/webdav/pkg/server"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the WebDAV server",
	Run: func(cmd *cobra.Command, args []string) {
		logging.Log.Info.Println("Creating user directories...")
		server.CreateUserDirectories()
		logging.Log.Info.Println("User directories created successfully")
		logging.Log.Info.Println("Hashing non-hashed passwords...")
		auth.HashPasswords()
		err := config.Write()
		if err != nil {
			logging.Log.Error.Fatalf("Error writing config: %s\n", err)
		}
		logging.Log.Info.Println("Passwords hashed successfully")
		logging.Log.Info.Println("Starting server...")
		err = server.StartWebdavServer()
		if err != nil {
			logging.Log.Error.Fatalf("Error starting server: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
