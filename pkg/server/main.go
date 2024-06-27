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

func StartWebdavServer() error {
	cfg := config.Get()
	dir := cfg.Content.Dir

	if !fs.PathExists(dir) {
		slog.Info("Creating data directory", "path", dir)
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	address := fmt.Sprintf("%s:%s", cfg.Network.Address, cfg.Network.Port)
	fileSystem := handler.NewWebdavFs(webdav.Dir(dir))
	if fileSystem == nil {
		slog.Error("Failed to create file system")
		return os.ErrInvalid
	}

	webdavSrv := handler.NewWebdavHandler(fileSystem, webdav.NewMemLS(), webdavLogger)
	http.Handle("/", auth.Middleware(webdavSrv))
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
