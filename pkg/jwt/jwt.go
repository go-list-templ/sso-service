package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-list-templ/sso-service/pkg/private_key"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secret = "secret"

const issuer = "sso"

var claims = jwt.MapClaims{
	"iss":        issuer,
	"sub":        "email",
	"iat":        time.Now().Unix(),
	"user_email": "email",
	"user_name":  "user",
}

var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

var ErrInvalidToken = errors.New("invalid token")

func keyFunc() jwt.Keyfunc {
	privateKey, _ := private_key.LoadPrivateKey("")

	return func(_ *jwt.Token) (interface{}, error) { return privateKey.PublicKey, nil }
}

func verifyAccessToken(accessToken string) error {
	token, err := jwt.Parse(accessToken,
		keyFunc(),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithIssuer(issuer),
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

func createAccessToken() (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"iss":        issuer,
		"sub":        "email",
		"iat":        now.Unix(),
		"exp":        now.Add(15 * time.Minute).Unix(),
		"user_email": "email",
		"user_name":  "name",
	}

	privateKey := ""

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func createRefreshToken() (string, error) {
	tokenID := uuid.New().String()
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":  issuer,
		"sub":  "email",
		"iat":  now.Unix(),
		"exp":  now.Add(7 * 24 * time.Hour).Unix(),
		"jti":  tokenID,
		"type": "refresh",
	}

	privateKey := ""

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	/* Если бы мы использовали opaqueToken, то мы бы хранили его hash */

	return signed, nil
}

func verifyRefreshToken(refreshToken string) error {
	token, err := jwt.Parse(refreshToken, keyFunc(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return fmt.Errorf("parse token failed: %w", err)
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return ErrInvalidToken
	}

	tokenID, ok := claims["jti"].(string)
	if !ok {
		return ErrInvalidToken
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return ErrInvalidToken
	}

	fmt.Println(tokenID, email)

	//mx.RLock()
	//_, exists := refreshTokens[tokenID]
	//mx.RUnlock()
	//
	//// Проверяем, не отозвали ли токен
	//if !exists {
	//	return ErrInvalidToken
	//}
	//
	//mx.Lock()
	//delete(refreshTokens, tokenID) // больше нельзя использовать этот refresh (одноразовый)
	//mx.Unlock()
	//
	//idx := slices.IndexFunc(usersDB, func(u user) bool {
	//	return strings.EqualFold(email, u.Email)
	//})
	//if idx == -1 {
	//	return ErrInvalidToken
	//}

	return nil
}
