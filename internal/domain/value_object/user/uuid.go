package user_vo

import (
	"errors"

	"github.com/google/uuid"
)

type UUID struct {
	value uuid.UUID
}

func NewUUID(value string) (*UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}
	return &UUID{value: parsed}, nil
}

func (u *UUID) String() string {
	return u.value.String()
}

func (u *UUID) Value() uuid.UUID {
	return u.value
}
