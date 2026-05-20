package dto

type PublicKeys struct {
	Keys []PublicKey
}

type PublicKey struct {
	Version string
	Key     string
}
