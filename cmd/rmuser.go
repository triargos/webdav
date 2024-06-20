/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/logging"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// rmuserCmd represents the rmuser command
var rmuserCmd = &cobra.Command{
	Use:   "rmuser",
	Short: "Remove a user from the webdav server configuration",
	Run: func(cmd *cobra.Command, args []string) {
		username := cmd.Flag("username").Value.String()
		cfg := config.Get()
		user, ok := (*cfg.Users)[username]
		if !ok {
			logging.Log.Error.Fatalf("User %s does not exist", username)
		}
		dirPath := filepath.Join(cfg.Content.Dir, user.Root)
		removeDirErr := os.RemoveAll(dirPath)
		if removeDirErr != nil {
			logging.Log.Error.Fatalf("Error removing user root directory: %s\n", removeDirErr)
		}
		config.RemoveUser(username)
		err := config.Write()
		if err != nil {
			logging.Log.Error.Fatalf("Error writing config: %s\n", err)
		}
		logging.Log.Info.Printf("Removed user %s was successful. Please restart the service for changes to take effect\n", username)

	},
}

func init() {
	rootCmd.AddCommand(rmuserCmd)

	rmuserCmd.Flags().StringP("username", "u", "", "Username of the user to remove")
}
