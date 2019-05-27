-- table schema for acks DB
CREATE TABLE IF NOT EXISTS acks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_email STRING,
    message STRING NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);