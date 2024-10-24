package handler

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
)

type WebDAVHandler struct {
	*webdav.Handler
}

func NewWebdavHandler(fs webdav.FileSystem, ls webdav.LockSystem, logger func(*http.Request, error)) *WebDAVHandler {
	return &WebDAVHandler{
		Handler: &webdav.Handler{
			FileSystem: fs,
			LockSystem: ls,
			Logger:     logger,
		},
	}
}

func (h *WebDAVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodHead {
		h.handleHead(w, r)
		return
	}
	h.Handler.ServeHTTP(w, r)

}

func (h *WebDAVHandler) handleHead(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	ctx := r.Context()
	f, err := h.FileSystem.OpenFile(ctx, reqPath, os.O_RDONLY, 0)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if fi.IsDir() {
		etag, err := findETag(ctx, h.FileSystem, h.LockSystem, reqPath, fi)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("ETag", etag)
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusOK)
		return
	}
	etag, err := findETag(ctx, h.FileSystem, h.LockSystem, reqPath, fi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("ETag", etag)

	http.ServeContent(w, r, reqPath, fi.ModTime(), f)
}

func findETag(ctx context.Context, fs webdav.FileSystem, ls webdav.LockSystem, name string, fi os.FileInfo) (string, error) {
	if do, ok := fi.(webdav.ETager); ok {
		etag, err := do.ETag(ctx)
		if !errors.Is(err, webdav.ErrNotImplemented) {
			return etag, err
		}
	}
	return fmt.Sprintf(`"%x%x"`, fi.ModTime().UnixNano(), fi.Size()), nil
}
