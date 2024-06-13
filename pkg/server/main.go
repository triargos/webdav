package server

import (
	"context"
	"fmt"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/logging"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartWebdavServer() error {
	if !fs.PathExists(config.Value.Content.Dir) {
		logging.Log.Info.Println("Creating data directory...")
		err := os.Mkdir(config.Value.Content.Dir, 0755)
		if err != nil {
			return err
		}
	}
	address := fmt.Sprintf("%s:%s", config.Value.Network.Address, config.Value.Network.Port)
	webdavSrv := &webdav.Handler{
		FileSystem: &WebdavFs{webdav.Dir(config.Value.Content.Dir)},
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			logging.Log.Info.Printf("%s %s: %s\n", r.Method, r.URL, err)
		},
	}

	http.Handle("/", AuthMiddleware(webdavSrv))
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		if !ok || !AuthenticateUser(username, password) {
			writer.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !CheckPermission(request.URL.Path, username) {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(request.Context(), "user", username)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
