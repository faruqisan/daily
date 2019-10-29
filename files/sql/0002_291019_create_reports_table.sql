create table reports
(
	id bigserial
		constraint reports_pk
			primary key,
	user_id bigserial not null,
	title varchar,
	detail varchar,
	created_at timestamp default now()
);

alter table users owner to dev;
