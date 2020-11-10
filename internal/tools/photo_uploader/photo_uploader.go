package photo_uploader

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type PhotoUploader struct {
}

func (pu *PhotoUploader) UploadPhoto(ctx echo.Context, formValue string, dirPath string) (string, error) {
	formFile, err := ctx.FormFile(formValue)
	if err != nil {
		return "", err
	}

	source, err := formFile.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = source.Close()
	}()
	tempBuffer := make([]byte, int64(math.Min(512, float64(formFile.Size))))
	_, err = source.Read(tempBuffer)
	if err != nil {
		return "", err
	}

	var fileExtension string
	imageType := http.DetectContentType(tempBuffer)
	switch imageType {
	case "image/jpg":
		fileExtension = "jpg"
		break
	case "image/jpeg":
		fileExtension = "jpeg"
		break
	case "image/png":
		fileExtension = "png"
		break
	case "image/webp":
		fileExtension = "webp"
		break
	default:
		return "", errors.New(fmt.Sprintf("Тип картинки %s не поддерживается", imageType))
	}

	hash := md5.Sum([]byte(formFile.Filename))
	fileName := hex.EncodeToString(hash[:]) + "." + fileExtension
	pathToNewFile := dirPath + fileName
	destination, err := os.OpenFile(pathToNewFile, os.O_WRONLY|os.O_CREATE, os.FileMode(int(0777)))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = destination.Close()
	}()

	if _, err = destination.Write(tempBuffer); err != nil {
		return "", err
	}

	if _, err = io.Copy(destination, source); err != nil {
		return "", err
	}

	return pathToNewFile, nil
}
