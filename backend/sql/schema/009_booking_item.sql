-- +goose Up
CREATE TABLE booking_item(
    booking_item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    service_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price_cents INTEGER NOT NULL,
    points_award INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_booking_item_service
        FOREIGN KEY (service_id)
        REFERENCES service (service_id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_booking_item_quantity
        CHECK (quantity > 0),

    CONSTRAINT chk_booking_item_unit_price_cents
        CHECK (unit_price_cents >= 0),

    CONSTRAINT chk_booking_item_points_award
        CHECK (points_award >= 0),

    -- Optional but recommended: prevents duplicate service lines per booking
    CONSTRAINT uq_booking_item_booking_service
        UNIQUE (booking_id, service_id)
);

-- +goose Down
DROP TABLE booking_item;
