package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"meguru-backend/internal/domain/entity"
	repository_interface "meguru-backend/internal/domain/repository"
	user_vo "meguru-backend/internal/domain/value_object/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository_interface.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (user_id, name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.UserID.Value().String(),
		user.Name.String(),
		user.Email.String(),
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, user_id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &entity.User{}
	var userID string
	var emailStr string
	var nameStr string

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&userID,
		&emailStr,
		&user.PasswordHash,
		&nameStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	userIDVO, err := user_vo.NewUuid(userID)
	if err != nil {
		return nil, err
	}
	user.UserID = *userIDVO

	emailVO, err := user_vo.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}
	user.Email = *emailVO

	nameVO, err := user_vo.NewUserName(nameStr)
	if err != nil {
		return nil, err
	}
	user.Name = *nameVO

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, user_id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &entity.User{}
	var userID string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &userID, &user.Email, &user.PasswordHash, &user.Name,
		&user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	userIDVO, err := user_vo.NewUuid(userID)
	if err != nil {
		return nil, err
	}
	user.UserID = *userIDVO

	return user, nil
}
