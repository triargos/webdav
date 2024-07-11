/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/environment"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/user"
	"log/slog"
	"os"
)

// rmuserCmd represents the rmuser command
var rmuserCmd = &cobra.Command{
	Use:   "rmuser",
	Short: "Remove a user from the webdav server configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configService := config.NewViperConfigService(environment.NewOsEnvironmentService())
		username := cmd.Flag("username").Value.String()
		fsService := fs.NewOsFileSystemService()
		userService := user.NewOsUserService(configService, fsService)
		removeUserErr := userService.RemoveUser(username)
		if removeUserErr != nil {
			slog.Error("failed to remove user", "error", removeUserErr.Error())
			os.Exit(1)
		}
		slog.Info("Removed user successfully. Please restart the service for changes to take effect", "username", username)
	},
}

func init() {
	rootCmd.AddCommand(rmuserCmd)

	rmuserCmd.Flags().StringP("username", "u", "", "Username of the user to remove")
}
