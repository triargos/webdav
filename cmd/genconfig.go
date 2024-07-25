package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/environment"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/helper"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var genconfigCmd = &cobra.Command{
	Use:   "genconfig",
	Short: "Generate a new configuration",
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetString("port")
		dir, _ := cmd.Flags().GetString("dir")
		subdirectories, _ := cmd.Flags().GetStringSlice("subdirectories")
		authType, _ := cmd.Flags().GetString("type")
		isValidAuthType := helper.ValidateAuthType(authType)
		if !isValidAuthType {
			slog.Error("Auth type must be either 'basic' or 'digest'")
			os.Exit(1)
		}
		configValue := config.Config{
			Network: config.NetworkConfig{
				Address: address,
				Port:    port,
			},
			Content: config.ContentConfig{
				Dir:            dir,
				SubDirectories: subdirectories,
			},
			Security: config.SecurityConfig{
				AuthType: authType,
			},
			Users: map[string]config.User{},
		}
		fileSystemHandler := fs.NewOsFileSystemService()
		envService := environment.NewOsEnvironmentService()
		configService := config.NewConfigService(envService, fileSystemHandler)
		configService.Set(&configValue)
		writeErr := configService.Write()
		if writeErr != nil {
			slog.Error("Failed to write configuration file", "error", writeErr.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(genconfigCmd)
	genconfigCmd.Flags().StringP("type", "t", "basic", "Authentication type")
	genconfigCmd.Flags().StringP("address", "a", "0.0.0.0", "Address to bind the server to")
	genconfigCmd.Flags().StringP("port", "p", "8080", "Port to bind the server to")
	genconfigCmd.Flags().StringP("dir", "d", "/var/webdav/data", "Directory to store files in")
	genconfigCmd.Flags().StringSliceP("subdirectories", "s", []string{}, "Subdirectories to create in the content directory")

}
