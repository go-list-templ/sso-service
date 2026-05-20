package dto

type PublicKeys struct {
	Keys []PublicKey
}

type PublicKey struct {
	Version string
	Key     string
}

type SignJWT struct {
	Version   string
	Signature string
}
