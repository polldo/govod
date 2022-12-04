CREATE TABLE IF NOT EXISTS orders
(
	order_id      UUID                        NOT NULL,
	user_id       UUID                        NOT NULL,
	provider_id   TEXT                        NOT NULL,
	/* TODO: Can possible status values be modeled? */
	status        TEXT                        NOT NULL,
	/* Amount is kinda redundant (it's a derived data) but not a problem since order items cannot be modified. */
	amount        FLOAT                       NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (order_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS order_items
(
	order_id      UUID                        NOT NULL,
	course_id     UUID                        NOT NULL,
	price         FLOAT                       NOT NULL,
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (order_id, course_id),
	FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
	FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);
