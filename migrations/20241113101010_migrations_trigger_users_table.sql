-- +goose Up
-- +goose StatementBegin
create table access_list
(
    role         varchar not null,
    endpoint     varchar not null,
    unique(role, endpoint)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table access_list;
-- +goose StatementEnd
