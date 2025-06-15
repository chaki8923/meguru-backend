package store_vo

import (
	"errors"
)

type Street string

func NewStreet(value string) (*Street, error) {
	if value == "" {
		return nil, errors.New("street cannot be empty")
	}
	s := Street(value)
	return &s, nil
}

func (s *Street) Value() string {
	return string(*s)
}
