package cmd

import (
	"github.com/spf13/cobra"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/environment"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/handler"
	"github.com/triargos/webdav/pkg/server"
	"github.com/triargos/webdav/pkg/user"
	"golang.org/x/net/webdav"
	"log/slog"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the WebDAV server",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Initializing webdav server...")
		configService := config.NewViperConfigService(environment.NewOsEnvironmentService())
		fsService := fs.NewOsFileSystemService()
		userService := user.NewOsUserService(configService, fsService)
		slog.Info("Creating all user directories...")
		createDirectoryErr := userService.InitializeDirectories()
		if createDirectoryErr != nil {
			slog.Error("Failed to create user directories", "error", createDirectoryErr.Error())
			os.Exit(1)
		}
		slog.Info("Hashing unhashed passwords...")
		hashPasswordsErr := userService.HashPasswords()
		if hashPasswordsErr != nil {
			slog.Error("Failed to hash passwords", "error", hashPasswordsErr.Error())
			os.Exit(1)
		}
		slog.Info("Starting webdav server...")
		authService := auth.New(userService)
		webdavFileSystem := handler.NewWebdavFs(webdav.Dir(configService.Get().Content.Dir), authService)
		if webdavFileSystem == nil {
			slog.Error("Failed to create webdav filesystem")
			os.Exit(1)
		}
		startServerErr := server.StartWebdavServer(server.StartWebdavServerContainer{
			ConfigService:    configService,
			WebdavFileSystem: webdavFileSystem,
			AuthService:      authService,
		})
		if startServerErr != nil {
			slog.Error("Failed to start webdav server", "error", startServerErr.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
