-- +goose Up
-- +goose StatementBegin
create table users
(
    id         serial primary key,
    name       varchar   not null,
    email      varchar   not null,
    role       varchar   not null,
    password   varchar   not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
