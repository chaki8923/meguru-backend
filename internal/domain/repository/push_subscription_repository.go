package repository

import (
	"context"

	"meguru-backend/internal/domain/entity"
)

type PushSubscriptionRepository interface {
	Create(ctx context.Context, subscription *entity.PushSubscription) error
	FindByEndpoint(ctx context.Context, endpoint string) (*entity.PushSubscription, error)
	FindAll(ctx context.Context) ([]*entity.PushSubscription, error)
	Delete(ctx context.Context, endpoint string) error
	Update(ctx context.Context, subscription *entity.PushSubscription) error
}