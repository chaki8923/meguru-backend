-- 統合されたマイグレーションのロールバック - すべてのテーブルを削除
-- 外部キー制約があるため、依存関係の逆順で削除

DROP TABLE IF EXISTS store_products;
DROP TABLE IF EXISTS store_email_verification_tokens;
DROP TABLE IF EXISTS flyer_items;
DROP TABLE IF EXISTS campaign_stores;
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS flyers;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS push_subscriptions;
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS users;

-- 拡張機能も削除（必要に応じて）
-- DROP EXTENSION IF EXISTS "uuid-ossp"; 