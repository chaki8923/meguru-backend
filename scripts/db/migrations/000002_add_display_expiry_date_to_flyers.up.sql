-- Add display_expiry_date column to flyers table
ALTER TABLE flyers ADD COLUMN display_expiry_date TIMESTAMPTZ;

-- Add index for efficient filtering by expiry date
CREATE INDEX idx_flyers_display_expiry_date ON flyers(display_expiry_date);

-- Add comment for the new column
COMMENT ON COLUMN flyers.display_expiry_date IS 'チラシの表示期限。この日時を過ぎると関連商品は表示されない';
