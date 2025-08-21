package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"meguru-backend/internal/domain/entity"
)

type PasswordResetTokenRepository interface {
	// Create はパスワードリセットトークンを新規作成します
	Create(ctx context.Context, token *entity.PasswordResetToken) error
	
	// FindByToken はトークンでパスワードリセットトークンを取得します
	FindByToken(ctx context.Context, token string) (*entity.PasswordResetToken, error)
	
	// UpdateUsed はトークンの使用状態を更新します
	UpdateUsed(ctx context.Context, id uuid.UUID, used bool) error
	
	// DeleteExpired は期限切れのトークンを削除します
	DeleteExpired(ctx context.Context, now time.Time) error
	
	// DeleteByUserID は指定ユーザーの未使用トークンを全て削除します（新規発行時の重複防止）
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
