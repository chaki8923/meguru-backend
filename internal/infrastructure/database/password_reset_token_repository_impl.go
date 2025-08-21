package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
)

type passwordResetTokenRepositoryImpl struct {
	db *sql.DB
}

func NewPasswordResetTokenRepository(db *sql.DB) repository.PasswordResetTokenRepository {
	return &passwordResetTokenRepositoryImpl{db: db}
}

func (r *passwordResetTokenRepositoryImpl) Create(ctx context.Context, token *entity.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, used, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		token.ID.String(), token.UserID.String(), token.Token, token.ExpiresAt,
		token.Used, token.CreatedAt, token.UpdatedAt)

	return err
}

func (r *passwordResetTokenRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.PasswordResetToken, error) {
	query := `
		SELECT prt.id, prt.user_id, prt.token, prt.expires_at, prt.used, prt.created_at, prt.updated_at,
		       u.user_id, u.email, u.password_hash, u.name, u.created_at, u.updated_at
		FROM password_reset_tokens prt
		JOIN users u ON prt.user_id = u.user_id
		WHERE prt.token = $1 AND prt.used = false AND prt.expires_at > $2`

	resetToken := &entity.PasswordResetToken{}
	user := &entity.User{}
	var tokenIDStr, userIDStr, userIDStr2 string

	err := r.db.QueryRowContext(ctx, query, token, time.Now()).Scan(
		&tokenIDStr, &userIDStr, &resetToken.Token, &resetToken.ExpiresAt,
		&resetToken.Used, &resetToken.CreatedAt, &resetToken.UpdatedAt,
		&userIDStr2, &user.Email, &user.PasswordHash, &user.Name,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// UUIDの変換
	tokenID, err := uuid.Parse(tokenIDStr)
	if err != nil {
		return nil, err
	}
	resetToken.ID = tokenID

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	resetToken.UserID = userID

	userID2, err := uuid.Parse(userIDStr2)
	if err != nil {
		return nil, err
	}
	user.ID = userID2

	resetToken.User = *user

	return resetToken, nil
}

func (r *passwordResetTokenRepositoryImpl) UpdateUsed(ctx context.Context, id uuid.UUID, used bool) error {
	query := `
		UPDATE password_reset_tokens 
		SET used = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, used, id.String())
	return err
}

func (r *passwordResetTokenRepositoryImpl) DeleteExpired(ctx context.Context, now time.Time) error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < $1`
	_, err := r.db.ExecContext(ctx, query, now)
	return err
}

func (r *passwordResetTokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1 AND used = false`
	_, err := r.db.ExecContext(ctx, query, userID.String())
	return err
}
