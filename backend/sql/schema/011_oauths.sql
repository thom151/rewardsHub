-- +goose Up
CREATE TABLE user_oauth_connections (
    connection_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider VARCHAR(50) NOT NULL,         -- google
    service VARCHAR(50) NOT NULL,          -- calendar
    provider_email VARCHAR(255),
    refresh_token TEXT NOT NULL,
    access_token TEXT,
    expires_at TIMESTAMPTZ,
    scopes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT uq_user_provider_service
        UNIQUE (user_id, provider, service)
);


-- +goose Down
DROP TABLE user_oauth_connections;
