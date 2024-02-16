package main

import (
	"context"
	"fmt"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type customFs struct {
	webdav.FileSystem
}

type contextKey string

const userContextKey contextKey = "user"

func (fs *customFs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	config, _ := readConfiguration()
	username, _ := ctx.Value(userContextKey).(string)
	nameParts := strings.Split(name, "/")
	//Check if we are admin
	if username == config.AdminUserName {
		return fs.FileSystem.OpenFile(ctx, name, flag, perm)
	}
	//Check if we are at the users root
	if strings.Contains(name, config.UsersRoot) {
		//Check if its trying to access a user folder
		if len(nameParts) >= 3 && nameParts[2] != username {
			return nil, os.ErrPermission
		}
	}
	//Otherwise, check if we are in another restricted directory
	isPermitted, err := CheckPathPermission(name, username)
	if err != nil || !isPermitted {
		return nil, os.ErrPermission
	}
	return fs.FileSystem.OpenFile(ctx, name, flag, perm)
}

func main() {
	configuration, err := readConfiguration()
	if err != nil {
		log.Fatalf("Error reading configuration file: %s", err.Error())
	}
	initStorage(configuration.DataPath)
	webdavSrv := &webdav.Handler{
		FileSystem: &customFs{FileSystem: webdav.Dir(configuration.DataPath)},
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				fmt.Printf("WebDAV %s: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				fmt.Printf("WebDAV %s: %s \n", r.Method, r.URL)
			}
		},
	}
	http.Handle("/", credentialsMiddleware(webdavSrv))
	go func() {
		log.Println("WebDAV server listening at port 8080")
		if err := http.ListenAndServe(":80", nil); err != nil {
			log.Fatalf("WebDAV server failed on HTTP interface: %s", err)
		}
	}()

	//HTTPS if possible
	_, errCert := os.Stat(configuration.CertPath)
	_, errKey := os.Stat(configuration.KeyPath)
	if errCert == nil || errKey == nil {
		go func() {
			log.Println("WebDAV server listening at port 8443")
			httpsListenError := http.ListenAndServeTLS(":443", configuration.CertPath, configuration.KeyPath, nil)
			if httpsListenError != nil {
				log.Fatal("WebDAV server failed on HTTPS interface: ", httpsListenError.Error())
			}
		}()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	time.Sleep(2 * time.Second) // Wait for 2 seconds to finish processing requests
	log.Println("Server gracefully stopped")
}

func credentialsMiddleware(next http.Handler) http.Handler {
	config, err := readConfiguration()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Executing authentication middleware")
		username, password, _ := request.BasicAuth()
		if err != nil || !verifyCredentials(username, password, config.Realm) {
			log.Println("Could not authenticate")
			writer.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(request.Context(), userContextKey, username)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
