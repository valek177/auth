-- +goose Up
create table users
(
    id         serial primary key,
    name       varchar   not null,
    email      varchar   not null,
    role       varchar   not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table users;
