-- name: CreateMitt :one
INSERT INTO mitts (
    author, content
) VALUES (
    @author, @content
)
RETURNING *;

-- name: GetMitt :one
SELECT * FROM mitts
WHERE id = @id
LIMIT 1;

-- name: GetAllUserMitts :many
SELECT * FROM mitts
WHERE author = @author
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: UpdateMitt :one
UPDATE mitts
SET
    content = @content,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteMitt :exec
DELETE FROM mitts
WHERE id = @id;
