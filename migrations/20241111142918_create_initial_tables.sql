-- +goose Up
-- +goose StatementBegin
CREATE TYPE transaction_source AS ENUM ('game', 'server', 'payment');

CREATE TYPE transaction_state AS ENUM ('win', 'lose');

CREATE TABLE users (
	id bigserial NOT NULL,
	created_at timestamp with time zone NULL,
	updated_at timestamp with time zone NULL,
	deleted_at timestamp with time zone NULL,
	name text NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE balances (
	id bigserial NOT NULL,
	created_at timestamp with time zone NULL,
	updated_at timestamp with time zone NULL,
	deleted_at timestamp with time zone NULL,
	user_id bigint NOT NULL,
	amount text NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transactions (
	id bigserial NOT NULL,
	uid text NOT NULL UNIQUE,
	created_at timestamp with time zone NULL,
	updated_at timestamp with time zone NULL,
	deleted_at timestamp with time zone NULL,
	user_id bigint NOT NULL,
	source transaction_source,
	t_state transaction_state NOT NULL,
	amount text NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;

DROP TABLE balances;

DROP TABLE users;

DROP TYPE transaction_state;

DROP TYPE transaction_source;

-- +goose StatementEnd