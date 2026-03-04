-- +goose Up
CREATE TABLE job (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    assigned_to_user_id UUID,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_job_booking
         FOREIGN KEY (booking_id)
        REFERENCES booking (booking_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_job_assigned_user
        FOREIGN KEY (assigned_to_user_id)
        REFERENCES users (user_id)
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE job;
