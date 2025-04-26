-- name: CreateUser :one
INSERT INTO users (login, name) VALUES (@login, @name)
RETURNING id;
