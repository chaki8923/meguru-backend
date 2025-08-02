package database

import (
	"context"
	"database/sql"
	"errors"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
)

type pushSubscriptionRepository struct {
	db *sql.DB
}

func NewPushSubscriptionRepository(db *sql.DB) repository.PushSubscriptionRepository {
	return &pushSubscriptionRepository{db: db}
}

func (r *pushSubscriptionRepository) Create(ctx context.Context, subscription *entity.PushSubscription) error {
	query := `INSERT INTO push_subscriptions (endpoint, p256dh, auth, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, subscription.Endpoint, subscription.P256dh, subscription.Auth, subscription.CreatedAt).Scan(&subscription.ID)
	return err
}

func (r *pushSubscriptionRepository) FindByEndpoint(ctx context.Context, endpoint string) (*entity.PushSubscription, error) {
	query := `SELECT id, endpoint, p256dh, auth, created_at FROM push_subscriptions WHERE endpoint = $1`
	var subscription entity.PushSubscription
	err := r.db.QueryRowContext(ctx, query, endpoint).Scan(&subscription.ID, &subscription.Endpoint, &subscription.P256dh, &subscription.Auth, &subscription.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *pushSubscriptionRepository) FindAll(ctx context.Context) ([]*entity.PushSubscription, error) {
	query := `SELECT id, endpoint, p256dh, auth, created_at FROM push_subscriptions`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*entity.PushSubscription
	for rows.Next() {
		var subscription entity.PushSubscription
		if err := rows.Scan(&subscription.ID, &subscription.Endpoint, &subscription.P256dh, &subscription.Auth, &subscription.CreatedAt); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &subscription)
	}
	return subscriptions, nil
}

func (r *pushSubscriptionRepository) Delete(ctx context.Context, endpoint string) error {
	query := `DELETE FROM push_subscriptions WHERE endpoint = $1`
	_, err := r.db.ExecContext(ctx, query, endpoint)
	return err
}

func (r *pushSubscriptionRepository) Update(ctx context.Context, subscription *entity.PushSubscription) error {
	query := `UPDATE push_subscriptions SET p256dh = $1, auth = $2, created_at = $3 WHERE endpoint = $4`
	_, err := r.db.ExecContext(ctx, query, subscription.P256dh, subscription.Auth, subscription.CreatedAt, subscription.Endpoint)
	return err
}
