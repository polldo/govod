CREATE TABLE IF NOT EXISTS videos
(
	video_id      UUID                        NOT NULL,
	course_id     UUID                        NOT NULL,
	index         INT                         NOT NULL,
	name          TEXT                        NOT NULL,
	description   TEXT                        NOT NULL,
	free          BOOLEAN                     NOT NULL,
	url           TEXT                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	version       INT                         NOT NULL DEFAULT 1,

	PRIMARY KEY (video_id),
	FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
	UNIQUE(course_id, index)
);
