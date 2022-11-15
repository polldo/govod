CREATE TABLE IF NOT EXISTS tokens
(
	hash        BYTEA        NOT NULL,
	user_id     UUID         NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	expiry      TIMESTAMP    NOT NULL DEFAULT NOW(),
	scope       TEXT         NOT NULL,

	PRIMARY KEY (hash)
);
