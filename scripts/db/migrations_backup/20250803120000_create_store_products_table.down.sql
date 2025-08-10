-- Drop indexes
DROP INDEX IF EXISTS idx_store_products_status;
DROP INDEX IF EXISTS idx_store_products_product_id;
DROP INDEX IF EXISTS idx_store_products_store_id;

-- Drop table
DROP TABLE IF EXISTS store_products; 