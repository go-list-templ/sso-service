package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-list-templ/sso-service/pkg/vault"
	"strings"
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
	vault  *vault.Vault
}

func NewToken(cfg *config.Config, l *zap.Logger, v *vault.Vault) *Token {
	return &Token{cfg, l, v}
}

func (t *Token) CreateAccess() (string, error) {
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

	unsignedToken, err := token.SigningString()
	if err != nil {
		return "", err
	}

	signature, err := t.signWithVault(unsignedToken)
	if err != nil {
		return "", fmt.Errorf("vault sign error: %w", err)
	}

	return fmt.Sprintf("%s.%s", unsignedToken, signature), nil
}

func (t *Token) signWithVault(unsignedToken string) (string, error) {
	path := "transit/sign/sso-service-keys"

	data := map[string]interface{}{
		"input":                base64.StdEncoding.EncodeToString([]byte(unsignedToken)),
		"marshaling_algorithm": "jws",
	}

	secret, err := t.vault.Logical().Write(path, data)
	if err != nil {
		return "", err
	}

	rawSignature := secret.Data["signature"].(string)

	parts := strings.Split(rawSignature, ":")
	finalSignature := parts[len(parts)-1]

	return finalSignature, nil
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
