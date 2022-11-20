CREATE TABLE IF NOT EXISTS videos
(
	video_id      UUID                        NOT NULL,
	course_id     UUID                        NOT NULL,
	index         INT                         NOT NULL,
	name          TEXT                        NOT NULL,
	description   TEXT                        NOT NULL,
	free          BOOLEAN                     NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (video_id),
	FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
	UNIQUE(course_id, index)
);

CREATE TABLE IF NOT EXISTS videos_url
(
	video_id      UUID                        NOT NULL,
	url           TEXT                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (video_id),
	FOREIGN KEY (video_id) REFERENCES videos(video_id) ON DELETE CASCADE
);
