CREATE TYPE user_lvl AS ENUM ('default', 'beginner', 'intermediate', 'advanced');

CREATE TABLE IF NOT EXISTS users
(
	user_id BIGINT PRIMARY KEY,
	lvl user_lvl NOT NULL DEFAULT 'default'
);