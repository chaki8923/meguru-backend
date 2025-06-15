package store_vo

import (
	"errors"
	"regexp"
)

type PhoneNumber string

func NewPhoneNumber(value string) (*PhoneNumber, error) {
	r := regexp.MustCompile(`^\+?[0-9\-]{7,15}$`)
	if !r.MatchString(value) {
		return nil, errors.New("invalid phone number format")
	}
	phone := PhoneNumber(value)
	return &phone, nil
}

func (p *PhoneNumber) Value() string {
	return string(*p)
}
