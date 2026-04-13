-- name: RollbackDelete :one
update files set deleted = 0 where id = ?

returning *;
