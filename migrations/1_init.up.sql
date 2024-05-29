CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL UNIQUE ,
    telegram_id bigint NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    firstname text DEFAULT '' NOT NULL,
    lastname text DEFAULT '' NOT NULL,
    is_admin bool DEFAULT false not null,
    created_at timestamp DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
	id uuid DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
	telegram_id bigint NOT NULL,
	words text[] NOT NULL,
	target text NOT NULL,
	words_hash text NOT NULL,
	hash text NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);
