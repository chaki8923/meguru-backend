-- Remove indexes first
DROP INDEX IF EXISTS idx_users_zipcode;
DROP INDEX IF EXISTS idx_users_prefecture;

-- Remove address fields from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS zipcode,
DROP COLUMN IF EXISTS prefecture,
DROP COLUMN IF EXISTS city,
DROP COLUMN IF EXISTS street;
