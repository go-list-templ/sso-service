package dto

type (
	UserCreateInput struct {
		Email    string
		Password string
	}

	UserCreateOutput struct {
		UserId string
	}

	UserVerifyCredInput struct {
		Email    string
		Password string
	}

	UserVerifyCredOutput struct {
		UserId string
	}
)
