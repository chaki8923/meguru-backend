package database

import (
	"context"
	"database/sql"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (user_id, email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var dbID int64
	err := r.db.QueryRowContext(ctx, query,
		user.ID.String(), user.Email, user.PasswordHash, user.Name,
		user.CreatedAt, user.UpdatedAt).Scan(&dbID)

	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT user_id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &entity.User{}
	var userIDStr string
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&userIDStr, &user.Email, &user.PasswordHash, &user.Name,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	user.ID = parsedID

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT user_id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE user_id = $1`

	user := &entity.User{}
	var userIDStr string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&userIDStr, &user.Email, &user.PasswordHash, &user.Name,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	user.ID = parsedID

	return user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2`

	_, err := r.db.ExecContext(ctx, query, passwordHash, userID.String())
	return err
}
