create table users
(
	id bigserial not null
		constraint users_pk
			primary key,
	email varchar,
	login_method varchar,
	created_at timestamp default now()
);

alter table users owner to dev;

create unique index users_email_uindex
	on users (email);

