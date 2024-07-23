package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/environment"
	"github.com/triargos/webdav/pkg/fs"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var genconfigCmd = &cobra.Command{
	Use:   "resetconfig",
	Short: "Reset the current config to the default values",
	Long:  "Resets the current configuration file to the defaults specified in the github repository",
	Run: func(cmd *cobra.Command, args []string) {
		configService := config.NewConfigService(environment.NewOsEnvironmentService(), fs.NewOsFileSystemService())
		slog.Info("resetting configuration file")
		resetConfigErr := configService.Reset()
		if resetConfigErr != nil {
			slog.Error("Failed to reset configuration file", "error", resetConfigErr)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(genconfigCmd)

}
