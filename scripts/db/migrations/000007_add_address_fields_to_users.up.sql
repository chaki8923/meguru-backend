-- Add address fields to users table
ALTER TABLE users 
ADD COLUMN zipcode VARCHAR(255),
ADD COLUMN prefecture VARCHAR(255),
ADD COLUMN city VARCHAR(255),
ADD COLUMN street VARCHAR(255);

-- Add indexes for address fields
CREATE INDEX idx_users_zipcode ON users(zipcode) WHERE zipcode IS NOT NULL;
CREATE INDEX idx_users_prefecture ON users(prefecture) WHERE prefecture IS NOT NULL;

-- Add comment
COMMENT ON COLUMN users.zipcode IS 'ユーザーの住所（郵便番号）';
COMMENT ON COLUMN users.prefecture IS 'ユーザーの住所（都道府県）';
COMMENT ON COLUMN users.city IS 'ユーザーの住所（市区町村）';
COMMENT ON COLUMN users.street IS 'ユーザーの住所（番地）';
