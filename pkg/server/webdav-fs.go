package server

import (
	"context"
	"golang.org/x/net/webdav"
	"log"
	"os"
)

type WebdavFs struct {
	webdav.FileSystem
}

func (fileSystem *WebdavFs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Println("Opening file:", name)
	return fileSystem.FileSystem.OpenFile(ctx, name, flag, perm)
}
