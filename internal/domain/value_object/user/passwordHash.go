package user_vo

import "errors"

type PasswordHash struct {
	value string
}

func NewPasswordHash(hash string) (*PasswordHash, error) {
	if hash == "" {
		return nil, errors.New("password hash cannot be empty")
	}

	if len(hash) != 60 {
		return nil, errors.New("password hash length is invalid")
	}
	return &PasswordHash{value: hash}, nil
}

func (p *PasswordHash) String() string {
	return p.value
}

func (p *PasswordHash) Bytes() []byte {
	return []byte(p.value)
}
