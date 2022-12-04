CREATE TABLE IF NOT EXISTS payments
(
	payment_id    UUID                        NOT NULL,
	order_id      UUID                        NOT NULL,
	/* TODO: Consider switching to UUID instead. */
	provider_id   TEXT                        NOT NULL,
	/* TODO: Can possible status values be modeled? */
	status        TEXT                        NOT NULL,
	amount        FLOAT                       NOT NULL,
	/* currency      TEXT                        NOT NULL  DEFAULT "EUR", */
	created_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at    TIMESTAMP                   NOT NULL DEFAULT NOW(),
	/* version       INT                         NOT NULL DEFAULT 1, */

	PRIMARY KEY (payment_id),
	FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);
