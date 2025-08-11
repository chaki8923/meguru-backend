-- Remove display_expiry_date column from flyers table
DROP INDEX IF EXISTS idx_flyers_display_expiry_date;
ALTER TABLE flyers DROP COLUMN IF EXISTS display_expiry_date;
