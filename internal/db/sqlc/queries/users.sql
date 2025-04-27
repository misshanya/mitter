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
