package dto

type (
	AuthInput struct {
		Email    string
		Password string
	}

	AuthOutput struct {
		AccessToken  string
		RefreshToken string
	}
)
