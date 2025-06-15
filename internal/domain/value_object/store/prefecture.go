package store_vo

import (
	"errors"
)

type Prefecture string

var validPrefectures = map[string]struct{}{
	"Tokyo": {}, "Osaka": {}, "Kyoto": {},
}

func NewPrefecture(value string) (*Prefecture, error) {
	if _, ok := validPrefectures[value]; !ok {
		return nil, errors.New("invalid prefecture")
	}
	p := Prefecture(value)
	return &p, nil
}

func (p *Prefecture) Value() string {
	return string(*p)
}
