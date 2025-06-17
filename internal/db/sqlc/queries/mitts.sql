-- name: CreateMitt :one
WITH inserted AS (
    INSERT INTO mitts (
        author, content
    ) VALUES (
        @author, @content
    )
    RETURNING *
)
SELECT i.*, us.name AS author_name, COUNT(ms.id) AS likes_count
FROM inserted AS i
LEFT JOIN mitts_likes AS ms ON ms.mitt_id = i.id
LEFT JOIN users AS us ON us.id = i.author
GROUP BY
    i.id,
    i.author,
    i.content,
    i.created_at,
    i.updated_at,
    us.name
LIMIT 1;

-- name: GetMitt :one
SELECT m.*, us.name AS author_name, COUNT(ms.id) AS likes_count
FROM mitts AS m
LEFT JOIN mitts_likes AS ms ON ms.mitt_id = m.id
LEFT JOIN users AS us ON us.id = m.author
WHERE m.id = @id
GROUP BY
    m.id,
    m.author,
    m.content,
    m.created_at,
    m.updated_at,
    us.name
LIMIT 1;

-- name: GetAllUserMitts :many
SELECT m.*, us.name AS author_name, COUNT(ms.id) AS likes_count
FROM mitts AS m
LEFT JOIN mitts_likes AS ms ON ms.mitt_id = m.id
LEFT JOIN users AS us ON us.id = m.author
WHERE m.author = @author
GROUP BY
    m.id,
    m.author,
    m.content,
    m.created_at,
    m.updated_at,
    us.name
ORDER BY m.created_at
LIMIT $1 OFFSET $2;

-- name: UpdateMitt :one
WITH updated AS (
    UPDATE mitts m
    SET
        content = @content,
        updated_at = NOW()
    WHERE m.id = @id
    RETURNING *
)
SELECT u.*, us.name AS author_name, COUNT(ms.id) AS likes_count
FROM updated AS u
LEFT JOIN mitts_likes AS ms ON ms.mitt_id = u.id
LEFT JOIN users AS us ON us.id = u.author
GROUP BY
    u.id,
    u.author,
    u.content,
    u.created_at,
    u.updated_at,
    us.name
LIMIT 1;

-- name: DeleteMitt :exec
DELETE FROM mitts
WHERE id = @id;


-- name: LikeMitt :exec
INSERT INTO mitts_likes (
    user_id, mitt_id
) VALUES (
    @user_id, @mitt_id
);

-- name: DeleteMittLike :exec
DELETE FROM mitts_likes
WHERE user_id = @user_id AND mitt_id = @mitt_id;


-- name: Feed :many
SELECT m.*, us.name AS author_name, COUNT(ms.id) AS likes_count
FROM mitts AS m
LEFT JOIN mitts_likes AS ms ON ms.mitt_id = m.id
LEFT JOIN users AS us ON us.id = m.author
GROUP BY
    m.id,
    m.author,
    m.content,
    m.created_at,
    m.updated_at,
    us.name
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
