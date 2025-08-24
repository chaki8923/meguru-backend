-- Create flyer_views table to track which users viewed which flyers
CREATE TABLE IF NOT EXISTS flyer_views (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    flyer_id UUID NOT NULL REFERENCES flyers(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Ensure one view record per user per flyer
    UNIQUE(flyer_id, user_id)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_flyer_views_flyer_id ON flyer_views(flyer_id);
CREATE INDEX IF NOT EXISTS idx_flyer_views_user_id ON flyer_views(user_id);
CREATE INDEX IF NOT EXISTS idx_flyer_views_viewed_at ON flyer_views(viewed_at);
