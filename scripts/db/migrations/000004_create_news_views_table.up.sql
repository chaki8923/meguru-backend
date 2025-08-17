-- Create news_views table to track RSS news reading history
CREATE TABLE news_views (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    news_url TEXT NOT NULL,
    news_title TEXT NOT NULL,
    news_id VARCHAR(255) NOT NULL, -- RSS feed内のID
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- 同じ店舗が同じニュースを短時間で重複閲覧することを防ぐためのユニーク制約
    UNIQUE(store_id, news_id)
);

-- Create indexes for performance
CREATE INDEX idx_news_views_store_id ON news_views(store_id);
CREATE INDEX idx_news_views_news_id ON news_views(news_id);
CREATE INDEX idx_news_views_viewed_at ON news_views(viewed_at);
CREATE INDEX idx_news_views_news_url ON news_views(news_url);

-- Create a composite index for spam check queries
CREATE INDEX idx_news_views_spam_check ON news_views(store_id, news_id, viewed_at);
