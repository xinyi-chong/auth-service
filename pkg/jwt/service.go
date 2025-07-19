package token

import (
	"auth-service/internal/config"
	"auth-service/pkg/logger"
	redisclient "auth-service/pkg/redis"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

var (
	secretKey     []byte
	configOnce    sync.Once
	accessExpiry  = time.Hour
	refreshExpiry = 30 * 24 * time.Hour
)

type (
	RefreshTokenClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	AccessTokenClaims struct {
		UserID   string  `json:"user_id"`
		Username *string `json:"username"`
		Email    *string `json:"email"`
		jwt.RegisteredClaims
	}
)

func Init(cfg *config.Config) error {
	var err error
	configOnce.Do(func() {
		if cfg.JWT.AccessDuration > 0 {
			accessExpiry = cfg.JWT.AccessDuration
		} else {
			logger.Warn("Invalid access duration, using default value: ", zap.Duration("accessExpiry", accessExpiry))
		}

		if cfg.JWT.RefreshDuration > 0 {
			refreshExpiry = cfg.JWT.RefreshDuration
		} else {
			logger.Warn("Invalid refresh duration, using default : ", zap.Duration("refreshExpiry", refreshExpiry))
		}

		key := os.Getenv("JWT_SECRET_KEY")
		if key == "" {
			err = errors.New("JWT secret key is not set in environment variables")
			return
		}
		secretKey = []byte(key)
	})
	return err
}

func GenerateAccessToken(userID string, username, email *string) (string, error) {
	accessClaims := AccessTokenClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiry)),
		},
	}
	return generateToken(accessClaims)
}

func GenerateRefreshToken(userID string) (string, error) {
	refreshClaims := RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry)),
		},
	}
	return generateToken(refreshClaims)
}

func ParseAccessToken(accessToken string) (*AccessTokenClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := parsedAccessToken.Claims.(*AccessTokenClaims)
	if !ok || !parsedAccessToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func ParseRefreshToken(refreshToken string) (*RefreshTokenClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := parsedRefreshToken.Claims.(*RefreshTokenClaims)
	if !ok || !parsedRefreshToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func InvalidateToken(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("empty_token")
	}

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsedToken.Valid {
		return errors.New("invalid token claims")
	}

	if claims.ExpiresAt == nil {
		return errors.New("missing expiration claim")
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return errors.New("token already expired")
	}

	hashedToken := secureHash(token)
	if err := redisclient.Set(ctx, getRedisBlacklistKey(hashedToken), "revoked", ttl); err != nil {
		return fmt.Errorf("storage failure: %w", err)
	}

	return nil
}

func IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	hashedToken := secureHash(token)
	return redisclient.Exists(ctx, getRedisBlacklistKey(hashedToken))
}

func generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return secretKey, nil
}

func secureHash(data string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func getRedisBlacklistKey(tokenHash string) string {
	return "blacklist:" + tokenHash
}
