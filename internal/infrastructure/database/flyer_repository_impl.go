package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type FlyerRepositoryImpl struct {
	db *sql.DB
}

func NewFlyerRepository(db *sql.DB) repository.FlyerRepository {
	return &FlyerRepositoryImpl{db: db}
}

func (r *FlyerRepositoryImpl) SaveFlyer(ctx context.Context, flyer *entity.Flyer) (*entity.Flyer, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	flyer.ID = uuid.New()
	flyer.CreatedAt = time.Now()
	flyer.UpdatedAt = flyer.CreatedAt

	_, err = tx.ExecContext(ctx, "INSERT INTO flyers (id, image_data, created_at, updated_at) VALUES ($1, $2, $3, $4)", flyer.ID, flyer.ImageData, flyer.CreatedAt, flyer.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into flyers: %w", err)
	}

	// TODO: Implement saving of store, campaign, products, and flyer_items from flyerData

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return flyer, nil
}
