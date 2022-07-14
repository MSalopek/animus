-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE DOMAIN bcrypt AS CHAR(60)
	CHECK (length(VALUE) IN (59, 60));

CREATE DOMAIN email AS TEXT
	CHECK (VALUE~'^[^@]+@[^@]+$');

CREATE DOMAIN domain_slug AS TEXT
	CHECK (VALUE~'^[a-z0-9\-]{1,32}$');

-- +goose StatementBegin
CREATE FUNCTION pg_soft_delete()
	RETURNS trigger AS $$
	DECLARE
		command text := ' SET deleted_at = current_timestamp WHERE id = $1';
	BEGIN
		EXECUTE 'UPDATE ' || TG_TABLE_NAME || command USING OLD.id;
		RETURN NULL;
	END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TYPE upload_stage AS ENUM (
	'storage', 'ipfs'
);

CREATE TYPE key_access_rights AS ENUM (
	'r',
	'rw',
	'rwd'
);

CREATE TYPE token_type AS ENUM (
	'register_email',
	'reset_pass'
);


CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	username VARCHAR(32) UNIQUE NOT NULL,
	firstname VARCHAR(32),
	lastname VARCHAR(32),
	email EMAIL UNIQUE NOT NULL,
	password CHAR(60) NOT NULL,
	max_keys INT NOT NULL default 5,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,
	active BOOLEAN NOT NULL DEFAULT false,

	CONSTRAINT users_email_uniq UNIQUE (email),
	CONSTRAINT users_email_valid CHECK (email ~ '^[^@]+@[^@]+$'),
	CONSTRAINT users_updated_at CHECK (updated_at >= created_at),
	CONSTRAINT users_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	),

	CONSTRAINT users_pass_bcrypt_check CHECK (LENGTH(password) IN (59, 60))
);

CREATE TRIGGER users_soft_delete
  BEFORE DELETE ON users
  FOR EACH ROW EXECUTE PROCEDURE pg_soft_delete();


CREATE TABLE tokens (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	token uuid UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
	type TOKEN_TYPE NOT NULL DEFAULT 'register_email'::token_type,
	
	valid_from TIMESTAMP NOT NULL DEFAULT now(),
	valid_to TIMESTAMP NOT NULL DEFAULT now() + interval '3 day',
	is_used BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX tokens_user_id_idx ON tokens(user_id);


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

CREATE TRIGGER storage_soft_delete
  BEFORE DELETE ON storage
  FOR EACH ROW EXECUTE PROCEDURE pg_soft_delete();

CREATE INDEX storage_dirs_idx ON storage(dir);
CREATE INDEX storage_user_id_idx ON storage(user_id);
CREATE INDEX storage_cid_idx ON storage(cid);

CREATE TABLE gateways (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	name VARCHAR(64) NOT NULL,
	slug VARCHAR(32) NOT NULL,
	public_id uuid DEFAULT uuid_generate_v4() NOT NULL
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
	valid_from TIMESTAMP NOT NULL DEFAULT now(),
	valid_to TIMESTAMP NOT NULL,
	config JSONB DEFAULT '{}' :: JSONB,

	CONSTRAINT subscriptions_valid_to_valid_from CHECK (valid_to >= valid_from),
	CONSTRAINT subscriptions_created_at_valid_from CHECK (created_at <= valid_from),
	CONSTRAINT subscriptions_created_at_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	)
);

CREATE TRIGGER subscriptions_soft_delete
  BEFORE DELETE ON subscriptions
  FOR EACH ROW EXECUTE PROCEDURE pg_soft_delete();

CREATE TABLE user_subscriptions (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	public_id uuid DEFAULT uuid_generate_v4() NOT NULL,
	subscription_id BIGINT REFERENCES subscriptions(id),
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,
	valid_from TIMESTAMP NOT NULL DEFAULT now(),
	valid_to TIMESTAMP NOT NULL,

	CONSTRAINT user_subscriptions_valid_to_valid_from CHECK (valid_to >= valid_from),
	CONSTRAINT user_subscriptions_created_at_valid_from CHECK (created_at <= valid_from),
	CONSTRAINT user_subscriptions_created_at_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	)
);

CREATE TRIGGER user_subscriptions_soft_delete
  BEFORE DELETE ON user_subscriptions
  FOR EACH ROW EXECUTE PROCEDURE pg_soft_delete();

-- keys can be assigned assigned as read, read-write or read-write-delete
CREATE TABLE keys (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT REFERENCES users(id),
	client_key varchar(32) UNIQUE NOT NULL,
	client_secret VARCHAR(64) NOT NULL,
	rights KEY_ACCESS_RIGHTS NOT NULL DEFAULT 'r',
	disabled BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	updated_at TIMESTAMP NOT NULL DEFAULT now(),
	deleted_at TIMESTAMP,
	valid_from TIMESTAMP NOT NULL DEFAULT now(),
	valid_to TIMESTAMP NOT NULL DEFAULT 'infinity',

	CONSTRAINT keys_valid_to_valid_from CHECK (valid_to >= valid_from),
	CONSTRAINT keys_created_at_valid_from CHECK (created_at <= valid_from),
	CONSTRAINT keys_updated_at CHECK (updated_at >= created_at),
	CONSTRAINT keys_created_at_deleted_at CHECK (
		deleted_at IS NULL
		OR deleted_at >= created_at
	)
);

-- +goose StatementBegin
CREATE TRIGGER keys_soft_delete
	BEFORE DELETE ON keys
	FOR EACH ROW EXECUTE PROCEDURE pg_soft_delete();

CREATE OR REPLACE FUNCTION check_add_key_permitted() RETURNS trigger AS
$$
DECLARE
    can_add boolean;
BEGIN
    WITH alloc_keys AS (
        SELECT u.id, u.max_keys, count(k.*) AS used_keys
        FROM users u
                 LEFT JOIN keys k ON k.user_id = u.id
        WHERE k.deleted_at IS NULL
          AND u.id = 1
          AND U.deleted_at IS NULL
        GROUP BY u.id, u.max_keys
    )
    SELECT CASE
               WHEN
                   alloc_keys.max_keys > alloc_keys.used_keys THEN true
               ELSE false
               END
    FROM alloc_keys
    INTO can_add;

    -- can happen if user is soft deleted
    IF can_add IS NULL THEN
        RAISE 'cannot add key';
    ELSIF NOT can_add THEN
        RAISE EXCEPTION 'maximum key allocation reached';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER before_key_insert
    BEFORE INSERT
    ON keys
EXECUTE PROCEDURE check_add_key_permitted();

-- +goose Down
DROP TABLE keys;
DROP TABLE user_subscriptions;
DROP TABLE subscriptions;
DROP TABLE gateways;
DROP TABLE storage;
DROP TABLE tokens;
DROP TABLE users;
DROP TYPE upload_stage;
DROP TYPE key_access_rights;
DROP TYPE bcrypt;
DROP TYPE domain_slug;
DROP TYPE email;
DROP TYPE token_type;
DROP FUNCTION check_add_key_permitted;
DROP FUNCTION pg_soft_delete;
