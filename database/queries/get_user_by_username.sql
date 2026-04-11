-- name: GetUserByUname :one
select * from users where username = ?;
