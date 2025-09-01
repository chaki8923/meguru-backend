package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
)

type flyerViewRepositoryImpl struct {
	db *sql.DB
}

func NewFlyerViewRepository(db *sql.DB) repository.FlyerViewRepository {
	return &flyerViewRepositoryImpl{db: db}
}

func (r *flyerViewRepositoryImpl) CreateOrIgnore(ctx context.Context, flyerView *entity.FlyerView) error {
	query := `
		INSERT INTO flyer_views (id, flyer_id, user_id, viewed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (flyer_id, user_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		flyerView.ID,
		flyerView.FlyerID.String(),
		flyerView.UserID,
		flyerView.ViewedAt,
		flyerView.CreatedAt,
		flyerView.UpdatedAt,
	)

	return err
}

func (r *flyerViewRepositoryImpl) GetViewCountByFlyerID(ctx context.Context, flyerID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM flyer_views WHERE flyer_id = $1`
	
	var count int64
	err := r.db.QueryRowContext(ctx, query, flyerID.String()).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *flyerViewRepositoryImpl) GetViewsByFlyerID(ctx context.Context, flyerID uuid.UUID, limit int, offset int) ([]*entity.FlyerView, error) {
	query := `
		SELECT fv.id, fv.flyer_id, fv.user_id, fv.viewed_at, fv.created_at, fv.updated_at,
		       u.user_id, u.email, u.name, u.created_at, u.updated_at
		FROM flyer_views fv
		JOIN users u ON fv.user_id = u.user_id
		WHERE fv.flyer_id = $1
		ORDER BY fv.viewed_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, flyerID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flyerViews []*entity.FlyerView
	for rows.Next() {
		flyerView := &entity.FlyerView{}
		user := &entity.User{}
		var flyerIDStr, userIDStr2 string

		err := rows.Scan(
			&flyerView.ID, &flyerIDStr, &flyerView.UserID, &flyerView.ViewedAt,
			&flyerView.CreatedAt, &flyerView.UpdatedAt,
			&userIDStr2, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// UUID変換
		flyerID, err := uuid.Parse(flyerIDStr)
		if err != nil {
			return nil, err
		}
		flyerView.FlyerID = flyerID

		userID2, err := uuid.Parse(userIDStr2)
		if err != nil {
			return nil, err
		}
		user.ID = userID2

		flyerView.User = *user
		flyerViews = append(flyerViews, flyerView)
	}

	return flyerViews, nil
}

func (r *flyerViewRepositoryImpl) HasUserViewedFlyer(ctx context.Context, flyerID uuid.UUID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM flyer_views WHERE flyer_id = $1 AND user_id = $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, flyerID.String(), userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	
	return exists, nil
}
