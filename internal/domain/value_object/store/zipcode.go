package store_vo

import (
	"errors"
	"regexp"
)

type Zipcode string

func NewZipcode(value string) (*Zipcode, error) {
	r := regexp.MustCompile(`^\d{3}-?\d{4}$`)
	if !r.MatchString(value) {
		return nil, errors.New("invalid zipcode format")
	}
	zipcode := Zipcode(value)
	return &zipcode, nil
}

func (z *Zipcode) Value() string {
	return string(*z)
}
