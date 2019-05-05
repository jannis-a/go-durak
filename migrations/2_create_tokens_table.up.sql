create table if not exists tokens
(
    id         serial primary key,
    user_id    integer                  not null references users (id) on delete cascade,
    token      text unique              not null,
    created_at timestamp with time zone not null default current_timestamp
);
