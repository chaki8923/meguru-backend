package repository

import (
	"context"
	"meguru-backend/internal/domain/entity"
)

type FlyerRepository interface {
	SaveFlyer(ctx context.Context, flyer *entity.Flyer) (*entity.Flyer, error)
}
