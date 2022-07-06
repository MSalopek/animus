-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE DOMAIN bcrypt AS CHAR(60)
-- 	CHECK (length(VALUE) IN (59, 60));

-- CREATE DOMAIN email AS TEXT
-- 	CHECK (VALUE~'^[^@]+@[^@]+$');

-- CREATE DOMAIN domain_slug AS TEXT
-- 	CHECK (VALUE~'^[a-z0-9\-]{1,32}$');

CREATE TYPE upload_stage AS ENUM (
	'storage', 'ipfs'
);

CREATE TYPE key_access_rights AS ENUM (
	'r',
	'rw',
	'rwd'
);

CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	username VARCHAR(32) UNIQUE NOT NULL,
	firstname VARCHAR(32),
	lastname VARCHAR(32),
	email VARCHAR(256) NOT NULL,
	password VARCHAR(60) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,

	CONSTRAINT users_email_uniq UNIQUE (email),
	CONSTRAINT users_email_valid CHECK (email ~ '^[^@]+@[^@]+$'),
	CONSTRAINT users_updated_at CHECK (updated_at >= created_at),
	CONSTRAINT users_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	),
	CONSTRAINT users_pass_hash CHECK (LENGTH(password) IN (59, 60))
);

CREATE TABLE storage (
	id BIGSERIAL PRIMARY KEY,
	cid VARCHAR(512),
	user_id BIGINT REFERENCES users(id),
	name VARCHAR(1024) NOT NULL,
	dir BOOLEAN NOT NULL DEFAULT false, -- true if the CID is used by a directory
	public BOOLEAN NOT NULL DEFAULT false,
	storage_bucket VARCHAR(1024),
	storage_key VARCHAR(1024),
	hash VARCHAR(1024),
	upload_stage UPLOAD_STAGE,
	pinned BOOLEAN NOT NULL DEFAULT false,
	metadata JSONB DEFAULT '{}'::JSONB,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,

	CONSTRAINT storage_updated_at CHECK (updated_at >= created_at),
	CONSTRAINT storage_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	)
);

CREATE INDEX storage_dirs_idx ON storage(dir);
CREATE INDEX storage_user_cid_idx ON storage(user_id, cid);

CREATE TABLE gateways (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	name VARCHAR(64) NOT NULL,
	slug VARCHAR(32) NOT NULL,
	public_id uuid DEFAULT uuid_generate_v4() NOT NULL,

	CONSTRAINT gateways_slug_valid CHECK (slug ~ '^[a-z0-9\-]{1,32}$')
);

CREATE TABLE subscriptions (
	id BIGSERIAL PRIMARY KEY,
	public_id uuid DEFAULT uuid_generate_v4() NOT NULL,
	name VARCHAR(64) NOT NULL,
	promotion BOOLEAN NOT NULL DEFAULT FALSE,
	-- if true it's free
	price NUMERIC NOT NULL DEFAULT 0,
	currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,
	valid_from TIMESTAMP NOT NULL,
	valid_to TIMESTAMP NOT NULL,
	config JSONB DEFAULT '{}' :: JSONB,

	CONSTRAINT subscriptions_valid_to_valid_from CHECK (valid_to >= valid_from),
	CONSTRAINT subscriptions_created_at_valid_from CHECK (created_at <= valid_from),
	CONSTRAINT subscriptions_created_at_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	)
);

CREATE TABLE user_subscriptions (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	public_id uuid DEFAULT uuid_generate_v4() NOT NULL,
	subscription_id BIGINT REFERENCES subscriptions(id),
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,
	valid_from TIMESTAMP NOT NULL,
	valid_to TIMESTAMP NOT NULL,

	CONSTRAINT user_subscriptions_valid_to_valid_from CHECK (valid_to >= valid_from),
	CONSTRAINT user_subscriptions_created_at_valid_from CHECK (created_at <= valid_from),
	CONSTRAINT user_subscriptions_created_at_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	)
);

-- keys can be assigned assigned as read, read-write or read-write-delete
CREATE TABLE keys (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	client_key varchar(32) NOT NULL,
	client_secret VARCHAR(32) NOT NULL,
	rights KEY_ACCESS_RIGHTS NOT NULL DEFAULT 'r',
	disabled BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,
	valid_from TIMESTAMP NOT NULL,
	valid_to TIMESTAMP NOT NULL,

	CONSTRAINT keys_valid_to_valid_from CHECK (valid_to >= valid_from),
	CONSTRAINT keys_created_at_valid_from CHECK (created_at <= valid_from),
	CONSTRAINT keys_updated_at CHECK (updated_at >= created_at),
	CONSTRAINT keys_created_at_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	)
);

-- +goose Down
DROP TABLE keys;
DROP TABLE user_subscriptions;
DROP TABLE subscriptions;
DROP TABLE gateways;
DROP TABLE storage;
DROP TABLE users;
DROP TYPE upload_stage;
DROP TYPE key_access_rights;
