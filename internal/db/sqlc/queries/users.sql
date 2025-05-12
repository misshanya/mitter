-- name: CreateUser :one
INSERT INTO users (login, name, password) VALUES (@login, @name, @hashedPassword)
RETURNING id;

-- name: GetUserByLogin :one
SELECT * FROM users WHERE login = @login
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = @id
LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = @id;

-- name: UpdateUser :exec
UPDATE users
SET
    name = COALESCE(sqlc.narg('name'), name)
WHERE id = @id;

-- name: UpdatePassword :exec
UPDATE users
SET
    password = @password
WHERE id = @id;

-- name: GetCurrentPasswordHash :one
SELECT password FROM users WHERE id = @id;

-- name: FollowUser :exec
INSERT INTO users_follows (
    follower_id, followee_id
) VALUES (
    @follower_id, @followee_id
);

-- name: UnfollowUser :exec
DELETE FROM users_follows
WHERE follower_id = @follower_id AND
      followee_id = @followee_id;

-- name: GetUserFollows :many
SELECT followee_id FROM users_follows
WHERE follower_id = @follower_id
LIMIT $1 OFFSET $2;

-- name: GetUserFollowers :many
SELECT follower_id FROM users_follows
WHERE followee_id = @followee_id
LIMIT $1 OFFSET $2;

-- name: GetUserFriends :many
SELECT uf1.followee_id
FROM users_follows uf1
JOIN users_follows uf2 ON uf1.follower_id = uf2.followee_id AND uf1.followee_id = uf2.follower_id
WHERE uf1.follower_id = @id
LIMIT $1 OFFSET $2;
