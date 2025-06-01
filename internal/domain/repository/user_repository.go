package repository

import (
	"context"

	"meguru-backend/internal/domain/entity"
	"github.com/google/uuid"
)
//このコメントアウトはpush前に絶対消すこと
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
} 