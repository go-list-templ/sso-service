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

const (
	Algorithm = "RS256"
	Type      = "RSA"
	Use       = "sig"
)

type JWKS struct{}

func New() *JWKS {
	return &JWKS{}
}

func (j *JWKS) FromPublicKeys(publicKeys dto.PublicKeys) (dto.JWKS, error) {
	result := dto.JWKS{
		Keys: make([]dto.JWK, 0, len(publicKeys.Keys)),
	}

	for _, pubKey := range publicKeys.Keys {
		block, _ := pem.Decode([]byte(pubKey.Key))
		if block == nil {
			return dto.JWKS{}, errors.New("parse pem for version: " + pubKey.Version)
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return dto.JWKS{}, err
		}

		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return dto.JWKS{}, errors.New("key is not an rsa for version: " + pubKey.Version)
		}

		eBytes := big.NewInt(int64(rsaPub.E)).Bytes()

		eBase64 := base64.RawURLEncoding.EncodeToString(eBytes)
		nBase64 := base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes())

		jwk := dto.JWK{
			Alg: Algorithm,
			E:   eBase64,
			Kid: pubKey.Version,
			Kty: Type,
			N:   nBase64,
			Use: Use,
		}

		result.Keys = append(result.Keys, jwk)
	}

	return result, nil
}
