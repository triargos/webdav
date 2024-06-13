package server

import (
	"fmt"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/handler"
	"github.com/triargos/webdav/pkg/logging"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func webdavLogger(req *http.Request, err error) {
	if err != nil {
		logging.Log.Error.Println(err)
	} else {
		logging.Log.Info.Println(req.Method, req.URL.Path)
	}

}

func StartWebdavServer() error {
	cfg := config.Get()
	dir := cfg.Content.Dir

	if !fs.PathExists(dir) {
		logging.Log.Info.Println("Creating data directory...")
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	address := fmt.Sprintf("%s:%s", cfg.Network.Address, cfg.Network.Port)
	fileSystem := handler.NewWebdavFs(webdav.Dir(dir))
	if fileSystem == nil {
		logging.Log.Error.Println("Failed to create file system")
		return os.ErrInvalid
	}

	webdavSrv := handler.NewWebdavHandler(fileSystem, webdav.NewMemLS(), webdavLogger)
	http.Handle("/", auth.Middleware(webdavSrv))
	go func() {
		logging.Log.Info.Printf("WebDAV server started on %s\n", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			logging.Log.Error.Printf("Error starting server: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Log.Info.Println("Server is shutting down...")
	time.Sleep(1 * time.Second)
	logging.Log.Info.Println("Server stopped")
	return nil
}
