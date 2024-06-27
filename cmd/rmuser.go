/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"log/slog"
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
			slog.Error("Failed to create user, does not exist", "username", username)
			return
		}
		dirPath := filepath.Join(cfg.Content.Dir, user.Root)
		removeDirErr := os.RemoveAll(dirPath)
		if removeDirErr != nil {
			slog.Error("Error removing user directory: %s", removeDirErr)
			os.Exit(1)
		}
		config.RemoveUser(username)
		err := config.Write()
		if err != nil {
			slog.Error("Failed to write config file:", "error", err.Error())
			os.Exit(1)
		}
		slog.Info("Removed user successfully. Please restart the service for changes to take effect", "username", username)

	},
}

func init() {
	rootCmd.AddCommand(rmuserCmd)

	rmuserCmd.Flags().StringP("username", "u", "", "Username of the user to remove")
}
