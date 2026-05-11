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

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrGenerateToken = errors.New("generate token")
)

type Token struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewToken(cfg *config.Config, l *zap.Logger) *Token {
	return &Token{cfg, l}
}

func (t *Token) Unsigned() (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"iss":        t.cfg.App.Name,
		"sub":        "email",
		"iat":        now.Unix(),
		"exp":        now.Add(DurationAccessToken).Unix(),
		"user_email": "email",
		"user_name":  "name",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SigningString()
}

func (t *Token) CreateAccess(unsignedToken, signature string) (string, error) {
	if unsignedToken == "" || signature == "" {
		return "", errors.New("invalid token components")
	}

	accessToken := unsignedToken + "." + signature

	parser := jwt.NewParser()
	_, _, err := parser.ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("generated token is invalid: %w", err)
	}

	return accessToken, nil
}

func (t *Token) CreateRefresh() (string, error) {
	b := make([]byte, LengthRefreshToken)
	_, err := rand.Read(b)
	if err != nil {
		t.logger.Error(ErrGenerateToken.Error(), zap.Error(err))

		return "", ErrGenerateToken
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (t *Token) keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (interface{}, error) { return "", nil }
}

func (t *Token) verifyAccess(accessToken string) error {
	token, err := jwt.Parse(accessToken,
		t.keyFunc(),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithIssuer(t.cfg.App.Name),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return fmt.Errorf("parse token failed: %w", err)
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrInvalidToken
	}

	userEmail, _ := claims["user_email"].(string)
	userName, _ := claims["user_name"].(string)

	fmt.Println(userEmail, userName)

	return nil
}
