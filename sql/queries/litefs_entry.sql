-- name: CreateEntry :exec
INSERT INTO litefs_entry (id, parent_id, name, modtime, content)
VALUES ($1, $2, $3, $4, $5);

-- name: LookupRootEntry :one
SELECT *
FROM litefs_entry
WHERE parent_id IS NULL AND name = $1;

-- name: LookupEntry :one
SELECT *
FROM litefs_entry
WHERE parent_id = $1 AND name = $2;

-- name: GetEntry :one
SELECT *
FROM litefs_entry
WHERE id = $1;

-- name: ListRootEntries :many
SELECT *
FROM litefs_entry
WHERE parent_id IS NULL;

-- name: ListEntries :many
SELECT *
FROM litefs_entry
WHERE parent_id = $1;

-- name: RenameEntry :exec
UPDATE litefs_entry
SET modtime = $2, name = $3
WHERE id = $1;

-- name: UpdateEntryBlob :exec
UPDATE litefs_entry
SET modtime = $2, content = $3
WHERE id = $1;

-- name: DeleteEntry :exec
DELETE FROM litefs_entry
WHERE id = $1;
