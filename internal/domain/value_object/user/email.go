package user_vo

import (
	"errors"
	"regexp"
)

type Email string

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func NewEmail(value string) (*Email, error) {
	if !emailRegex.MatchString(value) {
		return nil, errors.New("invalid email format")
	}
	e := Email(value)
	return &e, nil
}

func (e Email) String() string {
	return string(e)
}
