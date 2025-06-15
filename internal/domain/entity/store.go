package entity

import (
	"time"

	store_vo "meguru-backend/internal/domain/value_object/store"
)

type Store struct {
	ID           int64
	StoreID      store_vo.Uuid
	Name         string
	Email        store_vo.Email
	PasswordHash string
	PhoneNumber  store_vo.PhoneNumber
	Zipcode      store_vo.Zipcode
	Prefecture   store_vo.Prefecture
	City         store_vo.City
	Street       store_vo.Street
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewStore(
	id int64,
	storeID string,
	name string,
	email string,
	passwordHash string,
	phoneNumber string,
	zipcode string,
	prefecture string,
	city string,
	street string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Store, error) {
	uuidVO, err := store_vo.NewUuid(storeID)
	if err != nil {
		return nil, err
	}
	emailVO, err := store_vo.NewEmail(email)
	if err != nil {
		return nil, err
	}
	phoneVO, err := store_vo.NewPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	zipcodeVO, err := store_vo.NewZipcode(zipcode)
	if err != nil {
		return nil, err
	}
	prefectureVO, err := store_vo.NewPrefecture(prefecture)
	if err != nil {
		return nil, err
	}
	cityVO, err := store_vo.NewCity(city)
	if err != nil {
		return nil, err
	}
	streetVO, err := store_vo.NewStreet(street)
	if err != nil {
		return nil, err
	}

	return &Store{
		ID:           id,
		StoreID:      *uuidVO,
		Name:         name,
		Email:        *emailVO,
		PasswordHash: passwordHash,
		PhoneNumber:  *phoneVO,
		Zipcode:      *zipcodeVO,
		Prefecture:   *prefectureVO,
		City:         *cityVO,
		Street:       *streetVO,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
