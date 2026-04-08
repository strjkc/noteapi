-- name: CreateUser :one
insert into users(username, password, created_at, updated_at) values(?, ?, ?, ?)

returning *;
