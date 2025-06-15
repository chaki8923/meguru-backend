// internal/domain/query_service/user_query_service.go
package query_service

import (
	"context"

	"meguru-backend/internal/usecase/query_model"

	"github.com/google/uuid"
)

type UserQueryService interface {
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*query_model.User, error)
}
