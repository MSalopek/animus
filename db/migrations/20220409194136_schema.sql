-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE DOMAIN bcrypt AS CHAR(60)
-- 	CHECK (length(VALUE) IN (59, 60));

-- CREATE DOMAIN email AS TEXT
-- 	CHECK (VALUE~'^[^@]+@[^@]+$');

-- CREATE DOMAIN domain_slug AS TEXT
-- 	CHECK (VALUE~'^[a-z0-9\-]{1,32}$');

CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	username VARCHAR(32) UNIQUE NOT NULL,
	firstname VARCHAR(32),
	lastname VARCHAR(32),
	email VARCHAR(256) NOT NULL,
	password VARCHAR(60) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
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
	local BOOLEAN NOT NULL DEFAULT false,
	local_path VARCHAR(1024),
	hash VARCHAR(1024),
	uploaded BOOLEAN DEFAULT NOT NULL false,
	pinned BOOLEAN DEFAULT NOT NULL false,
	-- JSON data referencing any previous file versions
	-- check the versions column to fetch any previous versions if they stil exist
	-- TODO: create better versioning strategies
	versions JSONB DEFAULT '{}'::JSONB,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP,

	CONSTRAINT storage_updated_at
		CHECK (updated_at>=created_at),
	CONSTRAINT storage_deleted_at
		CHECK (deleted_at IS NULL OR deleted_at>=created_at),
	CONSTRAINT storage_local_files_check
		CHECK (local AND LENGTH(local_path) <> 0)
);

CREATE TABLE gateways (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	name VARCHAR(64) NOT NULL,
	slug VARCHAR(32) NOT NULL,
	public_id uuid DEFAULT uuid_generate_v4() NOT NULL,

	CONSTRAINT gateways_slug_valid
		CHECK (slug~'^[a-z0-9\-]{1,32}$')
);


-- +goose Down
DROP TABLE gateways;
DROP TABLE storage;
DROP TABLE users;
