package store_vo

import (
	"errors"

	"github.com/google/uuid"
)

type Uuid struct {
	value uuid.UUID
}

func NewUuid(value string) (*Uuid, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}
	return &Uuid{value: parsed}, nil
}

func (u *Uuid) String() string {
	return u.value.String()
}

func (u *Uuid) Value() uuid.UUID {
	return u.value
}
