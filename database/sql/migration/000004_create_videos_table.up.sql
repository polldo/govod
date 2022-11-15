CREATE TABLE IF NOT EXISTS videos
(
	course_id     UUID                        NOT NULL,
	video_index   INT                         NOT NULL,
	name          TEXT                        NOT NULL,
	description   TEXT                        NOT NULL,
	free          BOOLEAN                     NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (course_id, video_index),
	FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS video_url
(
	course_id     UUID                        NOT NULL,
	video_index   INT                         NOT NULL,
	url           TEXT                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (course_id, video_index),
	FOREIGN KEY (course_id, video_index) REFERENCES videos(course_id, video_index) ON DELETE CASCADE
);
