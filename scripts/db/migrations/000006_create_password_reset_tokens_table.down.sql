-- インデックス削除
DROP INDEX IF EXISTS idx_password_reset_tokens_expires_at;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;
DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;

-- テーブル削除
DROP TABLE IF EXISTS password_reset_tokens;
