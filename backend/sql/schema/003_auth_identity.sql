-- +goose Up
CREATE TABLE auth_identity(
    identity_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_subject VARCHAR(255),
    provider_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,

    CONSTRAINT unique_provider_identity
        UNIQUE(provider, provider_subject)
);

-- +goose Down
DROP TABLE auth_identity;

