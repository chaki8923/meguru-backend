-- 統合されたマイグレーションファイル - すべてのテーブルを作成
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users Table （最新の構造）
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Stores Table （すべての変更を統合）
CREATE TABLE stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    phone_number VARCHAR(255),
    zipcode VARCHAR(255),
    prefecture VARCHAR(255),
    city VARCHAR(255),
    street VARCHAR(255),
    email_verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Flyers Table
CREATE TABLE flyers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_data BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. Campaigns Table
CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flyer_id UUID REFERENCES flyers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Campaign-Stores Junction Table
CREATE TABLE campaign_stores (
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    PRIMARY KEY (campaign_id, store_id)
);

-- 6. Products Table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    category VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. Flyer-Items Table
CREATE TABLE flyer_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    price_excluding_tax INT,
    price_including_tax INT,
    unit VARCHAR(50),
    restriction_note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 8. Push Subscriptions Table
CREATE TABLE push_subscriptions (
    id SERIAL PRIMARY KEY,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 9. Store Email Verification Tokens Table
CREATE TABLE store_email_verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 10. Store Products Table
CREATE TABLE store_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    price INTEGER NOT NULL DEFAULT 0,
    quantity INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    status VARCHAR(50) NOT NULL DEFAULT '在庫あり',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- 同じ店舗で同じ商品を複数登録することを防ぐ
    UNIQUE(store_id, product_id)
);

-- インデックスの作成
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_stores_name ON stores(name);
CREATE INDEX idx_stores_email_verified_at ON stores(email_verified_at) WHERE email_verified_at IS NOT NULL;
CREATE INDEX idx_store_verification_tokens_store_id ON store_email_verification_tokens(store_id);
CREATE INDEX idx_store_verification_tokens_token ON store_email_verification_tokens(token);
CREATE INDEX idx_store_verification_tokens_expires_at ON store_email_verification_tokens(expires_at);
CREATE INDEX idx_store_products_store_id ON store_products(store_id);
CREATE INDEX idx_store_products_product_id ON store_products(product_id);
CREATE INDEX idx_store_products_status ON store_products(status);

-- テーブルコメント
COMMENT ON TABLE store_email_verification_tokens IS '店舗メール認証トークン管理テーブル';
COMMENT ON COLUMN store_email_verification_tokens.token IS 'メール認証用の一意なトークン';
COMMENT ON COLUMN store_email_verification_tokens.expires_at IS 'トークンの有効期限（通常24時間）'; 