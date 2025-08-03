-- メール認証フィールドの追加をロールバック
DROP INDEX IF EXISTS idx_stores_email_verified_at;
ALTER TABLE stores DROP COLUMN IF EXISTS email_verified_at; 