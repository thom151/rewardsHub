-- +goose Up
CREATE TABLE job_asset (
    asset_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    job_id UUID NOT NULL,
    created_by_user_id UUID,

    asset_type VARCHAR(50) NOT NULL,

    dropbox_path TEXT NOT NULL,
    dropbox_file_id TEXT,

    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,

    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    checksum_sha256 VARCHAR(64),

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_job_asset_job
        FOREIGN KEY (job_id)
        REFERENCES job(job_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_job_asset_created_by_user
        FOREIGN KEY (created_by_user_id)
        REFERENCES users(user_id)
        ON DELETE SET NULL,

    CONSTRAINT uq_job_asset_job_dropbox_path
        UNIQUE (job_id, dropbox_path)
);

-- +goose Down
DROP TABLE job_asset;
