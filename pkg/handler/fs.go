package handler

import (
	"context"
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/helper"
	"golang.org/x/net/webdav"
	"os"
)

type WebdavFs struct {
	webdav.FileSystem
	authService auth.Service
}

func NewWebdavFs(fs webdav.FileSystem, authService auth.Service) *WebdavFs {
	return &WebdavFs{
		FileSystem:  fs,
		authService: authService,
	}
}

func (filesystem *WebdavFs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	username, ok := helper.GetUsernameFromContext(ctx)
	if !ok || !filesystem.authService.HasPermission(name, username) {
		return nil, os.ErrPermission
	}
	return filesystem.FileSystem.OpenFile(ctx, name, flag, perm)
}

func (filesystem *WebdavFs) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	username, ok := helper.GetUsernameFromContext(ctx)
	if !ok || !filesystem.authService.HasPermission(name, username) {
		return nil, os.ErrPermission
	}
	return filesystem.FileSystem.Stat(ctx, name)
}

func (filesystem *WebdavFs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	username, ok := helper.GetUsernameFromContext(ctx)
	if !ok || !filesystem.authService.HasPermission(name, username) {
		return os.ErrPermission
	}
	return filesystem.FileSystem.Mkdir(ctx, name, perm)
}

func (filesystem *WebdavFs) RemoveAll(ctx context.Context, name string) error {
	username, ok := helper.GetUsernameFromContext(ctx)
	if !ok || !filesystem.authService.HasPermission(name, username) {
		return os.ErrPermission
	}
	return filesystem.FileSystem.RemoveAll(ctx, name)
}

func (filesystem *WebdavFs) Rename(ctx context.Context, oldName, newName string) error {
	username, ok := helper.GetUsernameFromContext(ctx)
	if !ok || !filesystem.authService.HasPermission(oldName, username) || !filesystem.authService.HasPermission(newName, username) {
		return os.ErrPermission
	}
	return filesystem.FileSystem.Rename(ctx, oldName, newName)
}
