CREATE TABLE IF NOT EXISTS carts
(
	user_id       UUID                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	version       INT                         NOT NULL DEFAULT 1,

	PRIMARY KEY (user_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS cart_items
(
	user_id       UUID                        NOT NULL,
	course_id     UUID                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (user_id, course_id),
	FOREIGN KEY (user_id) REFERENCES carts(user_id) ON DELETE CASCADE,
	FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);
