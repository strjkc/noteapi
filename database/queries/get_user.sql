-- name: GetUser :one
select * from users where username = ?;
