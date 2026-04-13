-- name: RemoveUser :one
delete from users where id = ?

returning *;
