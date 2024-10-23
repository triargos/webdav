package cmd

import (
	"github.com/spf13/cobra"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/cookie"
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
		fsService := fs.NewOsFileSystemService()
		configService := config.NewConfigService(environment.NewOsEnvironmentService(), fsService)

		userService := user.NewOsUserService(configService, fsService)
		slog.Info("Creating system and content directories...")
		contentDir := configService.Get().Content.Dir
		createDirectoriesErr := fsService.CreateDirectories(contentDir, 0755)
		if createDirectoriesErr != nil {
			slog.Error("Failed to create content directory", "error", createDirectoriesErr.Error())
			os.Exit(1)
		}
		for _, subirectory := range configService.Get().Content.SubDirectories {
			createSubDirectoryErr := fsService.CreateDirectories(contentDir+"/"+subirectory, 0755)
			if createSubDirectoryErr != nil {
				slog.Error("Failed to create subdirectory", "error", createSubDirectoryErr.Error())
			}
		}
		slog.Info("Creating all user directories...")
		createDirectoryErr := userService.InitializeDirectories()
		if createDirectoryErr != nil {
			slog.Error("Failed to create user directories", "error", createDirectoryErr.Error())
			os.Exit(1)
		}
		slog.Info("Checking if passwords should be hashed...")
		hashPasswordsErr := userService.HashPasswords()
		if hashPasswordsErr != nil {
			slog.Error("Failed to hash passwords", "error", hashPasswordsErr.Error())
			os.Exit(1)
		}

		slog.Info("Starting webdav server...")
		webdavFileSystem := handler.NewWebdavFs(webdav.Dir(configService.Get().Content.Dir))
		if webdavFileSystem == nil {
			slog.Error("Failed to create webdav filesystem")
			os.Exit(1)
		}
		digestAuthenticator := auth.NewDigestAuthenticator(userService)
		sslConfig := parseSSLConfig(cmd)
		cookieService := cookie.New()
		envService := environment.NewOsEnvironmentService()
		startServerErr := server.StartWebdavServer(server.StartWebdavServerContainer{
			ConfigService:       configService,
			WebdavFileSystem:    webdavFileSystem,
			UserService:         userService,
			FsService:           fsService,
			DigestAuthenticator: digestAuthenticator,
			SSLConfig:           sslConfig,
			CookieService:       cookieService,
			EnvService:          envService,
		})
		if startServerErr != nil {
			slog.Error("Failed to start webdav server", "error", startServerErr.Error())
			os.Exit(1)
		}
	},
}

func parseSSLConfig(cmd *cobra.Command) *server.SSLConfig {
	certFilePath, _ := cmd.Flags().GetString("certificate")
	keyFilePath, _ := cmd.Flags().GetString("key")
	if certFilePath != "" && keyFilePath != "" {
		return &server.SSLConfig{
			CertFilePath: certFilePath,
			KeyFilePath:  keyFilePath,
		}
	}
	return nil

}

func init() {
	startCmd.Flags().StringP("certificate", "c", "", "The certificate file path for HTTPS (optional)")
	startCmd.Flags().StringP("key", "k", "", "The key file path for HTTPS (optional)")

	rootCmd.AddCommand(startCmd)
}
