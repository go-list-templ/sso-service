package jwks

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"

	"github.com/go-list-templ/sso-service/internal/core/dto"
)

type JWKS struct{}

func New() *JWKS {
	return &JWKS{}
}

func pemToJWK(publicKeys dto.PublicKeys) (dto.JWK, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return dto.JWK{}, errors.New("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return dto.JWK{}, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return dto.JWK{}, errors.New("key is not an RSA public key")
	}

	eBytes := big.NewInt(int64(rsaPub.E)).Bytes()
	eBase64 := base64.RawURLEncoding.EncodeToString(eBytes)

	nBase64 := base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes())

	return dto.JWK{
		Alg: "RS256",
		E:   eBase64,
		Kid: version,
		Kty: "RSA",
		N:   nBase64,
		Use: "sig",
	}, nil
}
