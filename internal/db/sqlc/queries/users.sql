-- name: CreateUser :one
INSERT INTO users (login, name, password) VALUES (@login, @name, @hashedPassword)
RETURNING id;

-- name: GetUserByLogin :one
SELECT * FROM users WHERE login = @login
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = @id
LIMIT 1;
