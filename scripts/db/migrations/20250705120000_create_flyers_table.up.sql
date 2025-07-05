CREATE TABLE flyers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_data BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);