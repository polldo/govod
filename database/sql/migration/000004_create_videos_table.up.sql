CREATE TABLE IF NOT EXISTS videos
(
	video_id      UUID                        NOT NULL,
	course_id     UUID                        NOT NULL,
	index         INT                         NOT NULL,
	name          TEXT                        NOT NULL,
	description   TEXT                        NOT NULL,
	free          BOOLEAN                     NOT NULL,
	url           TEXT                        NOT NULL,
	image_url     TEXT                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	version       INT                         NOT NULL DEFAULT 1,

	PRIMARY KEY (video_id),
	FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
	UNIQUE(course_id, index)
);

CREATE TABLE IF NOT EXISTS videos_progress
(
	video_id      UUID                        NOT NULL,
	user_id       UUID                        NOT NULL,
	progress      INT                         NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	CHECK (progress BETWEEN 0 AND 100),
	PRIMARY KEY (video_id, user_id),
	FOREIGN KEY (video_id) REFERENCES videos(video_id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
