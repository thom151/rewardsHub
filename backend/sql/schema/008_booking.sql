-- +goose Up
CREATE TABLE booking(
    booking_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization(organization_id),
    requested_by_user_id UUID NOT NULL REFERENCES users(user_id),
    property_id UUID NOT NULL REFERENCES property(property_id),

    preferred_date DATE NOT NULL,
    schedule_start TIMESTAMPTZ NOT NULL,
    schedule_end TIMESTAMPTZ NOT NULL,

    status TEXT NOT NULL DEFAULT 'pending',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ,

    CHECK (schedule_end > schedule_start)
);

-- +goose Down
DROP TABLE booking;
