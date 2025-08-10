-- Create store_products table for store-specific product information
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

-- Indexes for better performance
CREATE INDEX idx_store_products_store_id ON store_products(store_id);
CREATE INDEX idx_store_products_product_id ON store_products(product_id);
CREATE INDEX idx_store_products_status ON store_products(status); 