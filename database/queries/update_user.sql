-- name: UpdateUser :one
update users set username = ?, password = ?, updated_at = ? where id = ?

returning *;
