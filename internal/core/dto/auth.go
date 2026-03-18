package dto

type (
	AuthInput struct {
		email    string
		password string
	}

	AuthOutput struct {
		token string
	}
)
