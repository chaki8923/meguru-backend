-- Remove index
DROP INDEX IF EXISTS idx_flyers_owner_store_id;

-- Remove foreign key constraint
ALTER TABLE flyers 
DROP CONSTRAINT IF EXISTS fk_flyers_owner_store;

-- Remove owner_store_id column
ALTER TABLE flyers 
DROP COLUMN IF EXISTS owner_store_id;
