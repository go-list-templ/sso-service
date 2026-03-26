package dto

import (
	"fmt"
	"net/mail"
)

const MinLenPass = 8

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

func (a *AuthInput) Validate() error {
	_, err := mail.ParseAddress(a.Email)
	if err != nil {
		return err
	}

	if len(a.Password) < MinLenPass {
		return fmt.Errorf("min lenght is %v chars", MinLenPass)
	}

	return nil
}
