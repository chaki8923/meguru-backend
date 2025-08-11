-- Drop indexes
DROP INDEX IF EXISTS idx_news_views_spam_check;
DROP INDEX IF EXISTS idx_news_views_news_url;
DROP INDEX IF EXISTS idx_news_views_viewed_at;
DROP INDEX IF EXISTS idx_news_views_news_id;
DROP INDEX IF EXISTS idx_news_views_store_id;

-- Drop news_views table
DROP TABLE IF EXISTS news_views;
