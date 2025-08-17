
-- 1. Tweets Table
CREATE TABLE tweets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    content VARCHAR(300) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Tweet Likes Table
CREATE TABLE tweet_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tweet_id UUID NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tweet_id, store_id)
);

-- Indexes
CREATE INDEX idx_tweets_store_id ON tweets(store_id);
CREATE INDEX idx_tweet_likes_tweet_id ON tweet_likes(tweet_id);
CREATE INDEX idx_tweet_likes_store_id ON tweet_likes(store_id);

-- Comments
COMMENT ON TABLE tweets IS '店舗のつぶやきを保存するテーブル';
COMMENT ON COLUMN tweets.content IS 'つぶやきの内容（300文字まで）';
COMMENT ON TABLE tweet_likes IS 'つぶやきへのいいねを記録するテーブル';
