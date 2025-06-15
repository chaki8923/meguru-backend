package user_vo

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type UserName string

func NewUserName(value string) (*UserName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, errors.New("user name cannot be empty")
	}
	length := utf8.RuneCountInString(value)
	if length < 2 || length > 50 {
		return nil, errors.New("user name must be between 2 and 50 characters")
	}

	n := UserName(value)
	return &n, nil
}

func (n UserName) String() string {
	return string(n)
}
