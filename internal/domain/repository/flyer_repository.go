package repository

import (
	"context"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/dto"
)

type FlyerRepository interface {
	SaveFlyer(ctx context.Context, flyer *entity.Flyer, flyerData *dto.FlyerData) (*entity.Flyer, error)
}
