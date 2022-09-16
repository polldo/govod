CREATE TABLE IF NOT EXISTS users
(
	id            UUID                        NOT NULL,
	name          TEXT                        NOT NULL,
	email         TEXT UNIQUE                 NOT NULL,
	role          TEXT                        NOT NULL,
	password_hash TEXT                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (id)
);
