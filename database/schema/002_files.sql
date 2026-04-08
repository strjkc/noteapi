-- +goose Up
create table files(id integer primary key autoincrement, name text not null, created_at text not null, updated_at text not null, user_id int not null references users(id) on delete cascade, unique(name, user_id));

-- +goose Down
drop table files;
