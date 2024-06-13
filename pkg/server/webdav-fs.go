package server

import (
	"context"
	"github.com/triargos/webdav/pkg/logging"
	"golang.org/x/net/webdav"
	"os"
)

type WebdavFs struct {
	webdav.FileSystem
}

func getUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value("user").(string)
	return username, ok
}

func (filesystem *WebdavFs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	username, ok := getUsernameFromContext(ctx)
	if !ok || !CheckPermission(name, username) {
		return nil, os.ErrPermission
	}
	logging.Log.Operation.Printf("User %s opened file %s", username, name)
	return filesystem.FileSystem.OpenFile(ctx, name, flag, perm)
}

func (filesystem *WebdavFs) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	username, ok := getUsernameFromContext(ctx)
	if !ok || !CheckPermission(name, username) {
		return nil, os.ErrPermission
	}
	logging.Log.Operation.Printf("User %s requested file info for %s", username, name)
	return filesystem.FileSystem.Stat(ctx, name)
}

func (filesystem *WebdavFs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	username, ok := getUsernameFromContext(ctx)
	if !ok || !CheckPermission(name, username) {
		return os.ErrPermission
	}
	logging.Log.Operation.Printf("User %s created directory %s", username, name)
	return filesystem.FileSystem.Mkdir(ctx, name, perm)
}

func (filesystem *WebdavFs) RemoveAll(ctx context.Context, name string) error {
	username, ok := getUsernameFromContext(ctx)
	if !ok || !CheckPermission(name, username) {
		return os.ErrPermission
	}
	logging.Log.Operation.Printf("User %s removed %s", username, name)
	return filesystem.FileSystem.RemoveAll(ctx, name)
}

func (filesystem *WebdavFs) Rename(ctx context.Context, oldName, newName string) error {
	username, ok := getUsernameFromContext(ctx)
	if !ok || !CheckPermission(oldName, username) || !CheckPermission(newName, username) {
		return os.ErrPermission
	}
	logging.Log.Operation.Printf("User %s renamed %s to %s", username, oldName, newName)
	return filesystem.FileSystem.Rename(ctx, oldName, newName)
}
