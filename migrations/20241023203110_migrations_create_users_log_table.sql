-- +goose Up
-- +goose StatementBegin
create table users_log
(
    id         serial primary key,
    user_id    int       not null,
    action     varchar   not null,
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users_log;
-- +goose StatementEnd
