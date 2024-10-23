package server

import (
	"fmt"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/cookie"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/handler"
	"github.com/triargos/webdav/pkg/user"
	"golang.org/x/net/webdav"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func webdavLogger(req *http.Request, err error) {
	if err != nil {
		slog.Error("REQ", "method", req.Method, "path", req.URL.Path, "error", err)
	} else {
		slog.Info("REQ", "method", req.Method, "path", req.URL.Path)
	}

}

type SSLConfig struct {
	CertFilePath string
	KeyFilePath  string
}

type StartWebdavServerContainer struct {
	ConfigService       config.Service
	UserService         user.Service
	DigestAuthenticator auth.DigestAuthenticator
	WebdavFileSystem    *handler.WebdavFs
	FsService           fs.Service
	SSLConfig           *SSLConfig
	CookieService       *cookie.Service
}

func StartWebdavServer(container StartWebdavServerContainer) error {
	configurationValue := container.ConfigService.Get()
	address := fmt.Sprintf("%s:%s", configurationValue.Network.Address, configurationValue.Network.Port)
	webdavSrv := handler.NewWebdavHandler(container.WebdavFileSystem, webdav.NewMemLS(), webdavLogger)
	authenticator := getAuthenticator(configurationValue, container.UserService)
	authMiddleware := auth.NewMiddleware(authenticator, container.CookieService)

	http.Handle("/", authMiddleware.Middleware(webdavSrv))
	go func() {
		if container.SSLConfig != nil {
			slog.Info("Starting the server using HTTPS...")
			if startErr := http.ListenAndServeTLS(":443", container.SSLConfig.CertFilePath, container.SSLConfig.KeyFilePath, nil); startErr != nil {
				slog.Error("failed to start HTTPS server", "error", startErr)
			}
		} else {
			slog.Info("Starting the server using HTTP...")
			if err := http.ListenAndServe(address, nil); err != nil {
				slog.Error("Failed to start server", "error", err)
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
	time.Sleep(1 * time.Second)
	slog.Info("Server stopped")
	return nil
}

func getAuthenticator(configurationValue *config.Config, userService user.Service) auth.Authenticator {
	if configurationValue.Security.AuthType == "digest" {
		return auth.NewDigestAuthenticator(userService)
	}
	return auth.NewBasicAuthenticator(userService)
}
