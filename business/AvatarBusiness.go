package business

import (
	"encoding/hex"
	"net/http"
	"crypto/md5"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/models"
)

var allowedType = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
}

func IsAllowedImageType(image []byte) (string, bool) {
	fileType := http.DetectContentType(image)
	extension, isAllowed := allowedType[fileType]

	return extension, isAllowed
}

func GetFileName(user *models.User, fileExtension string) string {
	randBytes := md5.Sum([]byte(fileExtension + user.Name))
	randString := hex.EncodeToString(randBytes[:])
	return randString + "." + fileExtension
}
