package dto

type (
	RefreshInput struct {
		RefreshToken string
	}

	AuthInput struct {
		Email    string
		Password string
	}

	AuthOutput struct {
		AccessToken  string
		RefreshToken string
	}
)
