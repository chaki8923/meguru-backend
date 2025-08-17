-- Add owner_store_id to flyers table to track which store created/owns the flyer
ALTER TABLE flyers 
ADD COLUMN owner_store_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';

-- Add foreign key constraint to stores table
ALTER TABLE flyers 
ADD CONSTRAINT fk_flyers_owner_store 
FOREIGN KEY (owner_store_id) REFERENCES stores(id) ON DELETE CASCADE;

-- Create index for performance
CREATE INDEX idx_flyers_owner_store_id ON flyers(owner_store_id);

-- Update existing flyers to set owner_store_id based on campaign_stores
-- This will set the first store found in campaign_stores as the owner
UPDATE flyers 
SET owner_store_id = (
    SELECT cs.store_id 
    FROM campaigns c 
    JOIN campaign_stores cs ON c.id = cs.campaign_id 
    WHERE c.flyer_id = flyers.id 
    LIMIT 1
)
WHERE owner_store_id = '00000000-0000-0000-0000-000000000000';

-- Remove the default value constraint after updating existing records
ALTER TABLE flyers 
ALTER COLUMN owner_store_id DROP DEFAULT;
