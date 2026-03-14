-- +goose Up
CREATE TABLE reward (
reward_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

code VARCHAR(50) NOT NULL UNIQUE,
name VARCHAR(255) NOT NULL,
description TEXT,

service_id UUID NOT NULL,

points_cost INTEGER NOT NULL CHECK (points_cost > 0),

is_active BOOLEAN NOT NULL DEFAULT TRUE,

created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

CONSTRAINT fk_reward_service
    FOREIGN KEY (service_id)
    REFERENCES service(service_id)
    ON DELETE SET NULL

);

-- +goose Down
DROP TABLE reward;
