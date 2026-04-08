-- +goose Up
alter table files add column deleted integer not null check(deleted in (0, 1));

-- +goose Down
alter table files drop column deleted;


