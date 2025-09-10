package auth

import (
	"auth-service/internal/shared/consts"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func isPasswordMatch(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func setRefreshTokenCookie(c *gin.Context, value string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		consts.CookieRefreshToken,
		value,
		int((time.Hour * 24 * 30).Seconds()),
		"/",
		"",
		true,
		true,
	)
}

func clearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie(
		consts.CookieRefreshToken,
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
}
