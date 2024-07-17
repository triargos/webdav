package server

import (
	"fmt"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/handler"
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

type StartWebdavServerContainer struct {
	ConfigService       config.Service
	AuthService         auth.Service
	DigestAuthenticator auth.DigestAuthenticator
	WebdavFileSystem    *handler.WebdavFs
	FsService           fs.Service
}

func StartWebdavServer(container StartWebdavServerContainer) error {
	configurationValue := container.ConfigService.Get()
	address := fmt.Sprintf("%s:%s", configurationValue.Network.Address, configurationValue.Network.Port)
	webdavSrv := handler.NewWebdavHandler(container.WebdavFileSystem, webdav.NewMemLS(), webdavLogger)
	authType := container.ConfigService.Get().Security.AuthType
	middleware := auth.BasicAuthMiddleware(container.AuthService)
	if authType == "digest" {
		middleware = auth.DigestAuthMiddleware(container.DigestAuthenticator, container.AuthService)
	}
	http.Handle("/", middleware(webdavSrv))
	go func() {
		slog.Info("Starting server", "address", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			slog.Error("Failed to start server", "error", err)
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
