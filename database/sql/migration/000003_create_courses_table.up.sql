CREATE TABLE IF NOT EXISTS courses
(
	course_id     UUID                        NOT NULL,
	name          TEXT                        NOT NULL,
	description   TEXT                        NOT NULL,
	price         FLOAT                       NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (course_id)
);
