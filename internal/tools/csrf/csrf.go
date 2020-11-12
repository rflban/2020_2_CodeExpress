package csrf

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
)

func NewCSRFToken(session *models.Session) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(session.ID))

	if err != nil {
		return "", err
	}

	hashedSession := hex.EncodeToString(hasher.Sum(nil))

	token := fmt.Sprintf("%s:%d", hashedSession, session.Expire.Unix())
	return token, nil
}

func ValidateCSRFToken(session *models.Session, token string) *ErrorResponse {
	tokenData := strings.Split(token, ":")
	csrfErr := errors.New("Bad csrf token received")

	if len(tokenData) != 2 {
		return NewErrorResponse(ErrBadRequest, csrfErr)
	}

	hasher := sha256.New()
	_, err := hasher.Write([]byte(session.ID))

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	hashedSession := hex.EncodeToString(hasher.Sum(nil))

	if hashedSession != tokenData[0] {
		return NewErrorResponse(ErrBadRequest, csrfErr)
	}

	expires, err := strconv.Atoi(tokenData[1])
	if err != nil {
		return NewErrorResponse(ErrBadRequest, csrfErr)
	}

	if int64(expires) < time.Now().Unix() {
		return NewErrorResponse(ErrBadRequest, csrfErr)
	}

	return nil
}
