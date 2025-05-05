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
WHERE follower_id = @follower_id;

-- name: GetUserFollowers :many
SELECT follower_id FROM users_follows
WHERE followee_id = @followee_id;
