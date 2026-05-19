package dto

type VaultPublicKey struct {
	Version string
	Key     string
}

type SignJWT struct {
	Version   string
	Signature string
}
