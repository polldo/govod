CREATE TABLE IF NOT EXISTS orders
(
	order_id      UUID                        NOT NULL,
	user_id       UUID                        NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	/* updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(), */

	PRIMARY KEY (order_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS order_items
(
	/* item_id       UUID                        NOT NULL, */
	order_id      UUID                        NOT NULL,
	course_id     UUID                        NOT NULL,
	price         FLOAT                       NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	/* updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(), */

	/* PRIMARY KEY (item_id), */
	/* UNIQUE(order_id, course_id) */
	PRIMARY KEY (order_id, course_id),
	FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
	FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);
