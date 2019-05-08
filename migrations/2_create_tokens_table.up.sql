create table if not exists tokens
(
	id         serial primary key,
	user_id    integer                  not null references users (id) on delete cascade,
	token      text unique              not null,
	login_at   timestamp with time zone not null default current_timestamp,
	login_ip   inet                     not null,
	refresh_at timestamp with time zone
);
