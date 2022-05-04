-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE DOMAIN bcrypt AS CHAR(60)
-- 	CHECK (length(VALUE) IN (59, 60));

-- CREATE DOMAIN email AS TEXT
-- 	CHECK (VALUE~'^[^@]+@[^@]+$');

-- CREATE DOMAIN domain_slug AS TEXT
-- 	CHECK (VALUE~'^[a-z0-9\-]{1,32}$');

CREATE TYPE upload_stage AS ENUM (
	"storage",
	"ipfs"
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
	
	CONSTRAINT users_email_valid
		CHECK (email~'^[^@]+@[^@]+$'),
	CONSTRAINT users_updated_at
		CHECK (updated_at>=created_at),
	CONSTRAINT users_deleted_at
		CHECK (deleted_at IS NULL OR deleted_at>=created_at),
	CONSTRAINT users_pass_has
		CHECK (length(password) IN (59, 60))
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
	-- JSON data referencing any previous file versions
	-- check the versions column to fetch any previous versions if they stil exist
	-- TODO: create better versioning strategies
	versions JSONB DEFAULT '{}'::JSONB,
	metadata JSONB DEFAULT '{}'::JSONB,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,

	CONSTRAINT storage_updated_at
		CHECK (updated_at>=created_at),
	CONSTRAINT storage_deleted_at
		CHECK (deleted_at IS NULL OR deleted_at>=created_at),
);

CREATE INDEX storage_dirs_idx ON storage(dir);
CREATE INDEX storage_user_cid_idx ON storage(user_id, cid);

CREATE TABLE gateways (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	name VARCHAR(64) NOT NULL,
	slug VARCHAR(32) NOT NULL,
	public_id uuid DEFAULT uuid_generate_v4() NOT NULL,

	CONSTRAINT gateways_slug_valid
		CHECK (slug~'^[a-z0-9\-]{1,32}$')
);

CREATE TABLE subscriptions (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	promotion BOOLEAN NOT NULL DEFAULT false, -- if true it's free
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	price NUMERIC NOT NULL DEFAULT 0,
	currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
	valid_from TIMESTAMP NOT NULL,
	valid_to TIMESTAMP NOT NULL,
	config JSONB DEFAULT '{}'::JSONB,

	CONSTRAINT subscriptions_updated_at
		CHECK (updated_at>=created_at),
	CONSTRAINT subscriptions_valid_to_valid_from
		CHECK (valid_to>=valid_from)
);

-- +goose Down
DROP TABLE gateways;
DROP TABLE storage;
DROP TABLE subscriptions;
DROP TABLE users;
