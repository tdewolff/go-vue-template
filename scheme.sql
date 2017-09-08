CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT,
	email TEXT
);

CREATE TABLE IF NOT EXISTS social_tokens (
	id INTEGER NOT NULL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	provider TEXT,
	access_token TEXT,
	refresh_token TEXT
);