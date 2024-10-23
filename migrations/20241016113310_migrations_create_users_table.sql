-- +goose Up
-- +goose StatementBegin
create table users
(
    id         serial primary key,
    name       varchar   not null,
    email      varchar   not null,
    role       varchar   not null,
    created_at timestamp not null default now(),
    super_pole varchar,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose StatementBegin
create table users_log
(
    id         serial primary key,
    user_id    int       not null,
    action     varchar   not null,
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose StatementBegin
create or replace function update_modified_column()
returns trigger as $$
BEGIN
NEW.updated_at = now();
NEW.super_pole = "new";
return NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
create trigger update_modified_time before update on users for each row execute procedure update_modified_column();
-- +goose StatementEnd

-- +goose Down
drop table users;
