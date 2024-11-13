-- +goose Up
-- +goose StatementBegin
create table roles_users_access
(
    id           serial primary key,
    role         int not null,
    endpoint     varchar   not null,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table roles_users_access;
-- +goose StatementEnd
