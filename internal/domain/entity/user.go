package entity

import (
	"time"

	user_vo "meguru-backend/internal/domain/value_object/user"
)

type User struct {
	ID           int64
	UserID       user_vo.Uuid
	Name         user_vo.UserName
	Email        user_vo.Email
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(
	id int64,
	userID string,
	name string,
	email string,
	passwordHash string,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {

	uuidVO, err := user_vo.NewUuid(userID)
	if err != nil {
		return nil, err
	}

	nameVO, err := user_vo.NewUserName(name)
	if err != nil {
		return nil, err
	}

	emailVO, err := user_vo.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		UserID:       *uuidVO,
		Name:         *nameVO,
		Email:        *emailVO,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
