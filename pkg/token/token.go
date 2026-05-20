package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/go-list-templ/sso-service/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const (
	LengthRefreshToken  = 64
	DurationAccessToken = 15 * time.Minute
)

var ErrGenerateToken = errors.New("generate token")

type Token struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewToken(cfg *config.Config, l *zap.Logger) *Token {
	return &Token{cfg, l}
}

func (t *Token) Unsigned(userId, keyVersion string, sessionCreatedAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"iss": t.cfg.App.ServiceName,
		"sub": userId,
		"iat": sessionCreatedAt.Unix(),
		"nbf": sessionCreatedAt.Unix(),
		"exp": sessionCreatedAt.Add(DurationAccessToken).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyVersion

	return token.SigningString()
}

func (t *Token) CreateAccess(unsignedToken, signature string) (string, error) {
	if unsignedToken == "" || signature == "" {
		return "", errors.New("invalid token components")
	}

	return fmt.Sprintf("%v.%v", unsignedToken, signature), nil
}

func (t *Token) CreateRefresh() (string, error) {
	b := make([]byte, LengthRefreshToken)
	_, err := rand.Read(b)
	if err != nil {
		t.logger.Error(ErrGenerateToken.Error(), zap.Error(err))

		return "", ErrGenerateToken
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
