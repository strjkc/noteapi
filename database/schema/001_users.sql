-- +goose Up
create table users(id integer primary key autoincrement not null, username text not null, password text not null, created_at text not null, updated_at text not null, unique(username));

-- +goose Down
drop table users;
