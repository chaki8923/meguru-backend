package database

import (
	"context"
	"database/sql"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"github.com/google/uuid"
)

type storeEmailVerificationTokenRepository struct {
	db *sql.DB
}

// NewStoreEmailVerificationTokenRepository 新しい店舗メール認証トークンリポジトリを作成
func NewStoreEmailVerificationTokenRepository(db *sql.DB) repository.StoreEmailVerificationTokenRepository {
	return &storeEmailVerificationTokenRepository{
		db: db,
	}
}

// Create トークンを作成
func (r *storeEmailVerificationTokenRepository) Create(ctx context.Context, storeID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO store_email_verification_tokens (id, store_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		uuid.New(), storeID, token, expiresAt, time.Now())

	return err
}

// FindByToken トークンで検索
func (r *storeEmailVerificationTokenRepository) FindByToken(ctx context.Context, token string) (*entity.StoreEmailVerificationToken, error) {
	query := `
		SELECT id, store_id, token, expires_at, created_at
		FROM store_email_verification_tokens
		WHERE token = $1`

	tokenEntity := &entity.StoreEmailVerificationToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&tokenEntity.ID, &tokenEntity.StoreID, &tokenEntity.Token, 
		&tokenEntity.ExpiresAt, &tokenEntity.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tokenEntity, nil
}

// Delete トークンを削除
func (r *storeEmailVerificationTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM store_email_verification_tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteByStoreID 指定した店舗のトークンをすべて削除
func (r *storeEmailVerificationTokenRepository) DeleteByStoreID(ctx context.Context, storeID uuid.UUID) error {
	query := `DELETE FROM store_email_verification_tokens WHERE store_id = $1`
	_, err := r.db.ExecContext(ctx, query, storeID)
	return err
}

// DeleteExpired 期限切れトークンを削除
func (r *storeEmailVerificationTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM store_email_verification_tokens WHERE expires_at < $1`
	_, err := r.db.ExecContext(ctx, query, time.Now())
	return err
} 