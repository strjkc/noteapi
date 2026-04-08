-- name: CreateFile :one
insert into files(name, created_at, updated_at, user_id) values(?, ?, ?, ?)

returning *;
