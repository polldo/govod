CREATE TABLE IF NOT EXISTS subscription_plans
(
	plan_id              UUID                        NOT NULL,
	stripe_id            TEXT                        NOT NULL,
	paypal_id            TEXT                        NOT NULL,
	name                 TEXT                        NOT NULL,
	price                INT                         NOT NULL,
	months_recurrence    INT                         NOT NULL,
	created_at           TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at           TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (plan_id)
);

CREATE TABLE IF NOT EXISTS subscriptions
(
	subscription_id      UUID                        NOT NULL,
	plan_id              UUID                        NOT NULL,
	user_id              UUID                        NOT NULL,
	provider             TEXT                        NOT NULL,
	provider_id          TEXT                        NOT NULL,
	/* TODO: Can possible status values be modeled? */
	status               TEXT                        NOT NULL,
	expiry               TIMESTAMP                   NOT NULL DEFAULT NOW(),
	created_at           TIMESTAMP                   NOT NULL DEFAULT NOW(),
	updated_at           TIMESTAMP                   NOT NULL DEFAULT NOW(),

	PRIMARY KEY (subscription_id),
	FOREIGN KEY (plan_id) REFERENCES subscription_plans(plan_id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	UNIQUE(provider_id)
);
