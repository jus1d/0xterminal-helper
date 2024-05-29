CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL UNIQUE ,
    telegram_id bigint NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    firstname text DEFAULT '' NOT NULL,
    lastname text DEFAULT '' NOT NULL,
    stage int DEFAULT 0 NOT NULL,
    is_admin bool DEFAULT false not null,
    created_at timestamp DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
	id uuid DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
	words text[] NOT NULL,
	target text NOT NULL,
	words_hash text NOT NULL,
	game_hash text NOT NULL,
	played_by bigint NOT NULL,
	played_at timestamp DEFAULT now() NOT NULL,
	FOREIGN KEY (played_by) REFERENCES users(telegram_id)
);
