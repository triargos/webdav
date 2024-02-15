package main

import (
	"fmt"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configuration, err := readConfiguration()
	if err != nil {
		log.Fatalf("Error reading configuration file: %s", err.Error())
	}
	initStorage(configuration.DataPath)
	webdavSrv := &webdav.Handler{
		FileSystem: webdav.Dir(configuration.DataPath),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				fmt.Printf("WebDAV %s: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				fmt.Printf("WebDAV %s: %s \n", r.Method, r.URL)
			}
		},
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		WebdavHandler(writer, request, webdavSrv)
	})
	go func() {
		log.Println("WebDAV server listening at port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("WebDAV server failed on HTTP interface: %s", err)
		}
	}()

	//HTTPS if possible
	_, errCert := os.Stat(configuration.CertPath)
	_, errKey := os.Stat(configuration.KeyPath)
	if errCert == nil || errKey == nil {
		go func() {
			log.Println("WebDAV server listening at port 8443")
			httpsListenError := http.ListenAndServeTLS(":8443", configuration.CertPath, configuration.KeyPath, nil)
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

func WebdavHandler(writer http.ResponseWriter, request *http.Request, handler *webdav.Handler) {
	config, err := readConfiguration()
	if err != nil {
		log.Println("Error reading configuration file")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	username, password, _ := request.BasicAuth()
	if !verifyCredentials(username, password, config.Realm) {
		writer.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	isPermitted := CheckPathPermission(request.URL.Path, username)
	if !isPermitted {
		http.Error(writer, "Forbidden", http.StatusForbidden)
		return
	}
	serveFile(writer, request, handler)
	return
}

func serveFile(writer http.ResponseWriter, request *http.Request, handler *webdav.Handler) {
	writer.Header().Set("Timeout", "99999999")
	handler.ServeHTTP(writer, request)
}
