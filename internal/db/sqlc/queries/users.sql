-- name: CreateUser :one
INSERT INTO users (login, name, password) VALUES (@login, @name, @hashedPassword)
RETURNING id;
