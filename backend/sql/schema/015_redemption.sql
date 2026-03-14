-- +goose Up
CREATE TABLE redemption (
    redemption_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organization_id UUID NOT NULL,
    requested_by_user_id UUID NOT NULL,
    reward_id UUID NOT NULL,
    applied_booking_id UUID NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    points_cost INTEGER NOT NULL CHECK (points_cost > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,

    CONSTRAINT fk_redemption_organization
        FOREIGN KEY (organization_id)
        REFERENCES organization(organization_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_redemption_requested_by_user
        FOREIGN KEY (requested_by_user_id)
        REFERENCES users(user_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_redemption_reward
        FOREIGN KEY (reward_id)
        REFERENCES reward(reward_id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_redemption_booking
        FOREIGN KEY (applied_booking_id)
        REFERENCES booking(booking_id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_redemption_status
        CHECK (status IN ('pending', 'approved', 'cancelled'))
);

-- +goose Down
DROP TABLE redemption;
