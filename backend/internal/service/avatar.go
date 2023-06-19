package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"webportfolio/internal"
)

func (s *service) CreateAvatar(ctx context.Context, req internal.UploadAvatarRequest) error {
	img, err := req.File.Open()
	if err != nil {
		return err
	}
	defer img.Close()

	// TODO: resize/convert?
	folder := filepath.Join("content", "avatars")
	os.MkdirAll(folder, os.ModePerm)
	userID := strconv.FormatInt(req.UserID, 10)

	dst, err := os.Create(filepath.Join(folder, userID))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, img); err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteAvatar(ctx context.Context, req internal.DeleteAvatarRequest) error {
	userID := strconv.FormatInt(req.UserID, 10)
	return os.Remove(filepath.Join("content", "avatars", userID))
}
