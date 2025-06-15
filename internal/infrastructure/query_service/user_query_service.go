package query_service

import (
	"context"
	"database/sql"

	"meguru-backend/internal/usecase/query_model"
	"meguru-backend/internal/usecase/query_service"

	"github.com/google/uuid"
)

type userQueryService struct {
	db *sql.DB
}

func NewUserQueryService(db *sql.DB) query_service.UserQueryService {
	return &userQueryService{db: db}
}

func (s *userQueryService) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*query_model.User, error) {
	query := `
		SELECT id, user_id, email, name, created_at
		FROM users
		WHERE user_id = $1
	`

	user := &query_model.User{}
	err := s.db.QueryRowContext(ctx, query, userID.String()).Scan(
		&user.ID, &user.UserID, &user.Email, &user.Name, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
