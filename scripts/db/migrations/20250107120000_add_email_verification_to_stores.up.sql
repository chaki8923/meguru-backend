-- 店舗テーブルにメール認証フィールドを追加
-- 既存の店舗に影響しないようにNULL許可で追加
ALTER TABLE stores ADD COLUMN email_verified_at TIMESTAMP WITH TIME ZONE;

-- インデックスを追加（認証済み店舗の検索を高速化）
CREATE INDEX idx_stores_email_verified_at ON stores(email_verified_at) WHERE email_verified_at IS NOT NULL; 