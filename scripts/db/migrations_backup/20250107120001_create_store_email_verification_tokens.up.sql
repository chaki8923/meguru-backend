-- 店舗メール認証トークンテーブルを作成
CREATE TABLE store_email_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- パフォーマンス向上のためのインデックス
CREATE INDEX idx_store_verification_tokens_store_id ON store_email_verification_tokens(store_id);
CREATE INDEX idx_store_verification_tokens_token ON store_email_verification_tokens(token);
CREATE INDEX idx_store_verification_tokens_expires_at ON store_email_verification_tokens(expires_at);

-- 期限切れトークンの自動削除用（オプション）
COMMENT ON TABLE store_email_verification_tokens IS '店舗メール認証トークン管理テーブル';
COMMENT ON COLUMN store_email_verification_tokens.token IS 'メール認証用の一意なトークン';
COMMENT ON COLUMN store_email_verification_tokens.expires_at IS 'トークンの有効期限（通常24時間）'; 