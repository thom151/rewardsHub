-- +goose Up
CREATE TABLE point_ledger (
    ledger_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organization_id UUID NOT NULL,
    user_id UUID,

    event_type VARCHAR(50) NOT NULL,
    event_id UUID,

    points_delta INTEGER NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_point_ledger_organization
        FOREIGN KEY (organization_id)
        REFERENCES organization(organization_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_point_ledger_user
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE point_ledger;
