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


-- name: LikeMitt :exec
INSERT INTO mitts_likes (
    user_id, mitt_id
) VALUES (
    @user_id, @mitt_id
);

-- name: IsMittLikedByUser :one
SELECT 1 FROM mitts_likes
WHERE user_id = @user_id AND mitt_id = @mitt_id;

-- name: DeleteMittLike :exec
DELETE FROM mitts_likes
WHERE user_id = @user_id AND mitt_id = @mitt_id;

-- name: GetMittLikesCount :one
SELECT COUNT(*) FROM mitts_likes
WHERE mitt_id = @mitt_id;
