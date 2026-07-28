package user

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	address string
}

func NewEmail(address string) (Email, error) {
	if address == "" {
		return Email{}, errors.New("email address is required")
	}
	if !emailRegex.MatchString(address) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{address: address}, nil
}

func (e Email) String() string {
	return e.address
}

func (e Email) IsZero() bool {
	return e.address == ""
}
