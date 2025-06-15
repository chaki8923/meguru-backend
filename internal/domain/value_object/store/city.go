package store_vo

import (
	"errors"
)

type City string

func NewCity(value string) (*City, error) {
	if value == "" {
		return nil, errors.New("city cannot be empty")
	}
	c := City(value)
	return &c, nil
}

func (c *City) Value() string {
	return string(*c)
}
