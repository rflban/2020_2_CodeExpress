package photo_uploader

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/labstack/echo/v4"
)

type PhotoUploader struct {
}

func (pu *PhotoUploader) UploadPhoto(ctx echo.Context, formValue string, dirPath string) (string, error) {
	var path string

	file, err := ctx.FormFile(formValue)

	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hash := md5.Sum([]byte(file.Filename))
	path = dirPath + hex.EncodeToString(hash[:])
	dst, err := os.Create(path)

	if err != nil {
		return "", err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return path, nil
}
