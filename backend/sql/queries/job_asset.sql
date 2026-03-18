-- name: CreateJobAsset :one
INSERT INTO job_asset (
job_id,
created_by_user_id,
asset_type,
dropbox_path,
dropbox_file_id,
file_name,
mime_type,
size_bytes,
checksum_sha256,
is_primary,
sort_order
)
VALUES (
$1,
$2,
$3,
$4,
$5,
$6,
$7,
$8,
$9,
$10,
$11
)
RETURNING *;

