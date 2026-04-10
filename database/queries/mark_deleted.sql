-- name: Deleted :exec
update files set deleted = 1 where id = ?;
