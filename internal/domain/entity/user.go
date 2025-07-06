package entity

import (
	"time"

	user_vo "meguru-backend/internal/domain/value_object/user"
)

type User struct {
	ID           int64
	UserID       user_vo.UUID
	Name         user_vo.UserName
	Email        user_vo.Email
	PasswordHash user_vo.PasswordHash
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(
	userID string,
	name string,
	email string,
	passwordHash string,
) (*User, error) {
	uuidVO, err := user_vo.NewUUID(userID)
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

	passwordHashVO, err := user_vo.NewPasswordHash(passwordHash)
	if err != nil {
		return nil, err
	}

	return &User{
		UserID:       *uuidVO,
		Name:         *nameVO,
		Email:        *emailVO,
		PasswordHash: *passwordHashVO,
	}, nil
}
