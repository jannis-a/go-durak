create table if not exists users
(
	id        serial primary key,
	username  text unique,
	email     text unique,
	password  text,
	joined_at timestamp with time zone not null default current_timestamp
);
