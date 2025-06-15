package store_vo

import (
	"errors"
	"regexp"
)

type Email string

func NewEmail(value string) (*Email, error) {
	r := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !r.MatchString(value) {
		return nil, errors.New("invalid email format")
	}
	email := Email(value)
	return &email, nil
}

func (e Email) String() string {
	return string(e)
}

func (e *Email) Value() string {
	return string(*e)
}
