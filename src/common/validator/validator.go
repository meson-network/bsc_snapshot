package validator

import (
	"errors"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func ValidateEmail(email string) error {
	validate := validator.New()
	return validate.Var(email, "email")
}

//length 6-20
//must contain number and letter
//special character is optional
func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("password length must between 6 - 20")
	}
	var hasNumber, hasLetter bool
	for _, c := range password {
		if hasNumber && hasLetter {
			return nil
		}
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLetter(c):
			hasLetter = true
		}
	}
	return errors.New("password must contain both number and letter")
}
