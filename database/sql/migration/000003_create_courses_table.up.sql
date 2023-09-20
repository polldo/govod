CREATE TABLE IF NOT EXISTS courses
(
	course_id     UUID                        NOT NULL,
	name          TEXT                        NOT NULL,
	description   TEXT                        NOT NULL,
	price         INT                         NOT NULL,
	image_url     TEXT                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	version       INT                         NOT NULL DEFAULT 1,

	PRIMARY KEY (course_id)
);
