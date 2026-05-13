package dto

type (
	RefreshInput struct {
		AccessToken string
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
