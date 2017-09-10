CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS social_tokens (
	id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	provider TEXT NOT NULL,
	access_token TEXT NOT NULL,
	refresh_token TEXT,
	PRIMARY KEY (id, user_id, provider)
);