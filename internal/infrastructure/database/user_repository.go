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
		INSERT INTO users (user_id, email, password_hash, name, zipcode, prefecture, city, street, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	// 空文字列をNULLとして扱う
	var zipcode, prefecture, city, street interface{}
	if user.Zipcode == "" {
		zipcode = nil
	} else {
		zipcode = user.Zipcode
	}
	if user.Prefecture == "" {
		prefecture = nil
	} else {
		prefecture = user.Prefecture
	}
	if user.City == "" {
		city = nil
	} else {
		city = user.City
	}
	if user.Street == "" {
		street = nil
	} else {
		street = user.Street
	}

	var dbID int64
	err := r.db.QueryRowContext(ctx, query,
		user.ID.String(), user.Email, user.PasswordHash, user.Name,
		zipcode, prefecture, city, street,
		user.CreatedAt, user.UpdatedAt).Scan(&dbID)

	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT user_id, email, password_hash, name, zipcode, prefecture, city, street, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &entity.User{}
	var userIDStr string
	var zipcode, prefecture, city, street sql.NullString
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&userIDStr, &user.Email, &user.PasswordHash, &user.Name,
		&zipcode, &prefecture, &city, &street,
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

	// NullString を string に変換
	user.Zipcode = zipcode.String
	user.Prefecture = prefecture.String
	user.City = city.String
	user.Street = street.String

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT user_id, email, password_hash, name, zipcode, prefecture, city, street, created_at, updated_at
		FROM users
		WHERE user_id = $1`

	user := &entity.User{}
	var userIDStr string
	var zipcode, prefecture, city, street sql.NullString
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&userIDStr, &user.Email, &user.PasswordHash, &user.Name,
		&zipcode, &prefecture, &city, &street,
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

	// NullString を string に変換
	user.Zipcode = zipcode.String
	user.Prefecture = prefecture.String
	user.City = city.String
	user.Street = street.String

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

func (r *userRepository) UpdateProfile(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET name = $1, zipcode = $2, prefecture = $3, city = $4, street = $5, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $6`

	// 空文字列をNULLとして扱う
	var zipcode, prefecture, city, street interface{}
	if user.Zipcode == "" {
		zipcode = nil
	} else {
		zipcode = user.Zipcode
	}
	if user.Prefecture == "" {
		prefecture = nil
	} else {
		prefecture = user.Prefecture
	}
	if user.City == "" {
		city = nil
	} else {
		city = user.City
	}
	if user.Street == "" {
		street = nil
	} else {
		street = user.Street
	}

	_, err := r.db.ExecContext(ctx, query, 
		user.Name, zipcode, prefecture, city, street, 
		user.ID.String())
	return err
} 
