--
-- PostgreSQL database dump
--

-- Dumped from database version 12.8 (Ubuntu 12.8-1.pgdg20.10+1)
-- Dumped by pg_dump version 12.8 (Ubuntu 12.8-1.pgdg20.10+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: bcrypt; Type: DOMAIN; Schema: public; Owner: animus
--

CREATE DOMAIN public.bcrypt AS character(60)
	CONSTRAINT bcrypt_check CHECK ((length(VALUE) = ANY (ARRAY[59, 60])));


ALTER DOMAIN public.bcrypt OWNER TO animus;

--
-- Name: domain_slug; Type: DOMAIN; Schema: public; Owner: animus
--

CREATE DOMAIN public.domain_slug AS text
	CONSTRAINT domain_slug_check CHECK ((VALUE ~ '^[a-z0-9\-]{1,32}$'::text));


ALTER DOMAIN public.domain_slug OWNER TO animus;

--
-- Name: email; Type: DOMAIN; Schema: public; Owner: animus
--

CREATE DOMAIN public.email AS text
	CONSTRAINT email_check CHECK ((VALUE ~ '^[^@]+@[^@]+$'::text));


ALTER DOMAIN public.email OWNER TO animus;

--
-- Name: key_access_rights; Type: TYPE; Schema: public; Owner: animus
--

CREATE TYPE public.key_access_rights AS ENUM (
    'r',
    'rw',
    'rwd'
);


ALTER TYPE public.key_access_rights OWNER TO animus;

--
-- Name: upload_stage; Type: TYPE; Schema: public; Owner: animus
--

CREATE TYPE public.upload_stage AS ENUM (
    'storage',
    'ipfs'
);


ALTER TYPE public.upload_stage OWNER TO animus;

--
-- Name: check_add_key_permitted(); Type: FUNCTION; Schema: public; Owner: animus
--

CREATE FUNCTION public.check_add_key_permitted() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.check_add_key_permitted() OWNER TO animus;

--
-- Name: pg_soft_delete(); Type: FUNCTION; Schema: public; Owner: animus
--

CREATE FUNCTION public.pg_soft_delete() RETURNS trigger
    LANGUAGE plpgsql
    AS $_$
	DECLARE
		command text := ' SET deleted_at = current_timestamp WHERE id = $1';
	BEGIN
		EXECUTE 'UPDATE ' || TG_TABLE_NAME || command USING OLD.id;
		RETURN NULL;
	END;
$_$;


ALTER FUNCTION public.pg_soft_delete() OWNER TO animus;

--
-- Name: soft_delete(); Type: FUNCTION; Schema: public; Owner: animus
--

CREATE FUNCTION public.soft_delete() RETURNS trigger
    LANGUAGE plpgsql
    AS $_$
    DECLARE
      command text := ' SET deleted_at = current_timestamp WHERE id = $1';
    BEGIN
      EXECUTE 'UPDATE ' || TG_TABLE_NAME || command USING OLD.id;
      RETURN NULL;
    END;
  $_$;


ALTER FUNCTION public.soft_delete() OWNER TO animus;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: gateways; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.gateways (
    id bigint NOT NULL,
    user_id bigint,
    name character varying(64) NOT NULL,
    slug character varying(32) NOT NULL,
    public_id uuid DEFAULT public.uuid_generate_v4() NOT NULL
);


ALTER TABLE public.gateways OWNER TO animus;

--
-- Name: gateways_id_seq; Type: SEQUENCE; Schema: public; Owner: animus
--

CREATE SEQUENCE public.gateways_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.gateways_id_seq OWNER TO animus;

--
-- Name: gateways_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: animus
--

ALTER SEQUENCE public.gateways_id_seq OWNED BY public.gateways.id;


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


ALTER TABLE public.goose_db_version OWNER TO animus;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: animus
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.goose_db_version_id_seq OWNER TO animus;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: animus
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: keys; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.keys (
    id bigint NOT NULL,
    user_id bigint,
    client_key character varying(32) UNIQUE NOT NULL,
    client_secret character varying(64) NOT NULL,
    rights public.key_access_rights DEFAULT 'r'::public.key_access_rights NOT NULL,
    disabled boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    valid_from timestamp without time zone NOT NULL,
    valid_to timestamp without time zone NOT NULL,
    CONSTRAINT keys_created_at_deleted_at CHECK (((deleted_at IS NULL) OR (deleted_at >= created_at))),
    CONSTRAINT keys_created_at_valid_from CHECK ((created_at <= valid_from)),
    CONSTRAINT keys_updated_at CHECK ((updated_at >= created_at)),
    CONSTRAINT keys_valid_to_valid_from CHECK ((valid_to >= valid_from))
);


ALTER TABLE public.keys OWNER TO animus;

--
-- Name: keys_id_seq; Type: SEQUENCE; Schema: public; Owner: animus
--

CREATE SEQUENCE public.keys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.keys_id_seq OWNER TO animus;

--
-- Name: keys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: animus
--

ALTER SEQUENCE public.keys_id_seq OWNED BY public.keys.id;


--
-- Name: storage; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.storage (
    id bigint NOT NULL,
    cid character varying(512),
    user_id bigint,
    name character varying(1024) NOT NULL,
    dir boolean DEFAULT false NOT NULL,
    public boolean DEFAULT false NOT NULL,
    storage_bucket character varying(1024),
    storage_key character varying(1024),
    hash character varying(1024),
    upload_stage public.upload_stage,
    pinned boolean DEFAULT false NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    CONSTRAINT storage_deleted_at CHECK (((deleted_at IS NULL) OR (deleted_at >= created_at))),
    CONSTRAINT storage_updated_at CHECK ((updated_at >= created_at))
);


ALTER TABLE public.storage OWNER TO animus;

--
-- Name: storage_id_seq; Type: SEQUENCE; Schema: public; Owner: animus
--

CREATE SEQUENCE public.storage_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.storage_id_seq OWNER TO animus;

--
-- Name: storage_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: animus
--

ALTER SEQUENCE public.storage_id_seq OWNED BY public.storage.id;


--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.subscriptions (
    id bigint NOT NULL,
    public_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(64) NOT NULL,
    promotion boolean DEFAULT false NOT NULL,
    price numeric DEFAULT 0 NOT NULL,
    currency character varying(3) DEFAULT 'EUR'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    valid_from timestamp without time zone NOT NULL,
    valid_to timestamp without time zone NOT NULL,
    config jsonb DEFAULT '{}'::jsonb,
    CONSTRAINT subscriptions_created_at_deleted_at CHECK (((deleted_at IS NULL) OR (deleted_at >= created_at))),
    CONSTRAINT subscriptions_created_at_valid_from CHECK ((created_at <= valid_from)),
    CONSTRAINT subscriptions_valid_to_valid_from CHECK ((valid_to >= valid_from))
);


ALTER TABLE public.subscriptions OWNER TO animus;

--
-- Name: subscriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: animus
--

CREATE SEQUENCE public.subscriptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.subscriptions_id_seq OWNER TO animus;

--
-- Name: subscriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: animus
--

ALTER SEQUENCE public.subscriptions_id_seq OWNED BY public.subscriptions.id;


--
-- Name: user_subscriptions; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.user_subscriptions (
    id bigint NOT NULL,
    user_id bigint,
    public_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    subscription_id bigint,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    valid_from timestamp without time zone NOT NULL,
    valid_to timestamp without time zone NOT NULL,
    CONSTRAINT user_subscriptions_created_at_deleted_at CHECK (((deleted_at IS NULL) OR (deleted_at >= created_at))),
    CONSTRAINT user_subscriptions_created_at_valid_from CHECK ((created_at <= valid_from)),
    CONSTRAINT user_subscriptions_valid_to_valid_from CHECK ((valid_to >= valid_from))
);


ALTER TABLE public.user_subscriptions OWNER TO animus;

--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: animus
--

CREATE SEQUENCE public.user_subscriptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_subscriptions_id_seq OWNER TO animus;

--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: animus
--

ALTER SEQUENCE public.user_subscriptions_id_seq OWNED BY public.user_subscriptions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying(32) NOT NULL,
    firstname character varying(32),
    lastname character varying(32),
    email public.email NOT NULL,
    password public.bcrypt NOT NULL,
    max_keys integer DEFAULT 5 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    CONSTRAINT users_deleted_at CHECK (((deleted_at IS NULL) OR (deleted_at >= created_at))),
    CONSTRAINT users_email_valid CHECK (((email)::text ~ '^[^@]+@[^@]+$'::text)),
    CONSTRAINT users_pass_hash CHECK ((length((password)::bpchar) = ANY (ARRAY[59, 60]))),
    CONSTRAINT users_updated_at CHECK ((updated_at >= created_at))
);


ALTER TABLE public.users OWNER TO animus;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: animus
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO animus;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: animus
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: gateways id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.gateways ALTER COLUMN id SET DEFAULT nextval('public.gateways_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: keys id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.keys ALTER COLUMN id SET DEFAULT nextval('public.keys_id_seq'::regclass);


--
-- Name: storage id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.storage ALTER COLUMN id SET DEFAULT nextval('public.storage_id_seq'::regclass);


--
-- Name: subscriptions id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.subscriptions ALTER COLUMN id SET DEFAULT nextval('public.subscriptions_id_seq'::regclass);


--
-- Name: user_subscriptions id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.user_subscriptions ALTER COLUMN id SET DEFAULT nextval('public.user_subscriptions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: gateways; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.gateways (id, user_id, name, slug, public_id) FROM stdin;
\.


--
-- Data for Name: goose_db_version; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.goose_db_version (id, version_id, is_applied, tstamp) FROM stdin;
1	0	t	2022-04-09 20:10:30.384308
12	20220409194136	t	2022-07-07 16:53:56.346862
\.


--
-- Data for Name: keys; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.keys (id, user_id, client_key, client_secret, rights, disabled, created_at, updated_at, deleted_at, valid_from, valid_to) FROM stdin;
1	1	aEdcerBC7N1Y82N14bVEo7KWIiM0Ntje	dnkHJySctQFTm39TXy0kY3Z8awUBITK1deMwsdpMkzi06Dc7ZqxAeyNYVO9m7uZP	rwd	f	2022-07-07 16:54:01.959992	2022-07-07 16:54:01.959992	\N	2022-07-07 16:54:01.959992	infinity
\.


--
-- Data for Name: storage; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.storage (id, cid, user_id, name, dir, public, storage_bucket, storage_key, hash, upload_stage, pinned, metadata, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: subscriptions; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.subscriptions (id, public_id, name, promotion, price, currency, created_at, deleted_at, valid_from, valid_to, config) FROM stdin;
1	00000000-0000-0000-0000-000000000000	Free Plan	f	0	EUR	2022-07-07 16:54:01.956942	\N	2022-07-07 16:54:01.956942	infinity	{}
\.


--
-- Data for Name: user_subscriptions; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.user_subscriptions (id, user_id, public_id, subscription_id, created_at, deleted_at, valid_from, valid_to) FROM stdin;
1	1	00000000-0000-0000-0000-000000000000	1	2022-07-07 16:54:01.958178	\N	2022-07-07 16:54:01.958178	infinity
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.users (id, username, firstname, lastname, email, password, max_keys, created_at, updated_at, deleted_at) FROM stdin;
1	admin	Animus	Administrator	admin@example.com	$2a$12$IRWQnDUmZ.OjvTuLKaBNte09IcwrvtcPni1G4rBYVpZLW0WJOuPnC	5	2022-04-13 00:00:00	2022-04-13 00:00:00	\N
\.


--
-- Name: gateways_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.gateways_id_seq', 1, false);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 12, true);


--
-- Name: keys_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.keys_id_seq', 1, true);


--
-- Name: storage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.storage_id_seq', 1, false);


--
-- Name: subscriptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.subscriptions_id_seq', 1, true);


--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.user_subscriptions_id_seq', 1, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- Name: gateways gateways_pkey; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.gateways
    ADD CONSTRAINT gateways_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: keys keys_pkey; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.keys
    ADD CONSTRAINT keys_pkey PRIMARY KEY (id);


--
-- Name: storage storage_pkey; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.storage
    ADD CONSTRAINT storage_pkey PRIMARY KEY (id);


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (id);


--
-- Name: user_subscriptions user_subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_uniq; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_uniq UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: storage_dirs_idx; Type: INDEX; Schema: public; Owner: animus
--

CREATE INDEX storage_dirs_idx ON public.storage USING btree (dir);


--
-- Name: storage_user_cid_idx; Type: INDEX; Schema: public; Owner: animus
--

CREATE INDEX storage_user_cid_idx ON public.storage USING btree (user_id, cid);


--
-- Name: keys before_key_insert; Type: TRIGGER; Schema: public; Owner: animus
--

CREATE TRIGGER before_key_insert BEFORE INSERT ON public.keys FOR EACH STATEMENT EXECUTE FUNCTION public.check_add_key_permitted();


--
-- Name: keys keys_soft_delete; Type: TRIGGER; Schema: public; Owner: animus
--

CREATE TRIGGER keys_soft_delete BEFORE DELETE ON public.keys FOR EACH ROW EXECUTE FUNCTION public.pg_soft_delete();


--
-- Name: storage storage_soft_delete; Type: TRIGGER; Schema: public; Owner: animus
--

CREATE TRIGGER storage_soft_delete BEFORE DELETE ON public.storage FOR EACH ROW EXECUTE FUNCTION public.pg_soft_delete();


--
-- Name: subscriptions subscriptions_soft_delete; Type: TRIGGER; Schema: public; Owner: animus
--

CREATE TRIGGER subscriptions_soft_delete BEFORE DELETE ON public.subscriptions FOR EACH ROW EXECUTE FUNCTION public.pg_soft_delete();


--
-- Name: user_subscriptions user_subscriptions_soft_delete; Type: TRIGGER; Schema: public; Owner: animus
--

CREATE TRIGGER user_subscriptions_soft_delete BEFORE DELETE ON public.user_subscriptions FOR EACH ROW EXECUTE FUNCTION public.pg_soft_delete();


--
-- Name: users users_soft_delete; Type: TRIGGER; Schema: public; Owner: animus
--

CREATE TRIGGER users_soft_delete BEFORE DELETE ON public.users FOR EACH ROW EXECUTE FUNCTION public.pg_soft_delete();


--
-- Name: gateways gateways_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.gateways
    ADD CONSTRAINT gateways_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: keys keys_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.keys
    ADD CONSTRAINT keys_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: storage storage_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.storage
    ADD CONSTRAINT storage_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_subscriptions user_subscriptions_subscription_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES public.subscriptions(id);


--
-- Name: user_subscriptions user_subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

