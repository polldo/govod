CREATE TABLE IF NOT EXISTS users
(
	user_id       UUID                        NOT NULL,
	name          TEXT                        NOT NULL,
	email         TEXT UNIQUE                 NOT NULL,
	role          TEXT                        NOT NULL,
	active        BOOLEAN                     NOT NULL,
	password_hash TEXT                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	version       INT                         NOT NULL DEFAULT 1,

	PRIMARY KEY (user_id)
);
