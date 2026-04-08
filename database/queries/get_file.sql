-- name: GetFile :one
select * from files where name = ? and user_id = ?;
