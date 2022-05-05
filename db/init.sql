--
-- PostgreSQL database dump
--

-- Dumped from database version 14.2
-- Dumped by pg_dump version 14.2

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


CREATE TYPE public.upload_stage AS ENUM (
	'storage',
	'ipfs'
);

--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


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
    public_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    CONSTRAINT gateways_slug_valid CHECK (((slug)::text ~ '^[a-z0-9\-]{1,32}$'::text))
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
-- Name: storage; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.storage (
    id bigint NOT NULL,
    cid character varying(512),
    user_id bigint,
    name character varying(1024) NOT NULL,
    dir boolean DEFAULT false NOT NULL,
    public boolean DEFAULT false NOT NULL,
	storage_bucket VARCHAR(1024),
	storage_key VARCHAR(1024),
	hash VARCHAR(1024),
	upload_stage public.upload_stage,
    pinned boolean DEFAULT false NOT NULL,
    versions jsonb DEFAULT '{}'::jsonb,
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
    user_id bigint,
    promotion boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    price numeric DEFAULT 0 NOT NULL,
    currency character varying(3) DEFAULT 'EUR'::character varying NOT NULL,
    valid_from timestamp without time zone NOT NULL,
    valid_to timestamp without time zone NOT NULL,
    config jsonb DEFAULT '{}'::jsonb,
    CONSTRAINT subscriptions_updated_at CHECK ((updated_at >= created_at)),
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
-- Name: users; Type: TABLE; Schema: public; Owner: animus
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying(32) NOT NULL,
    firstname character varying(32),
    lastname character varying(32),
    email character varying(256) NOT NULL,
    password character varying(60) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone,
    CONSTRAINT users_deleted_at CHECK (((deleted_at IS NULL) OR (deleted_at >= created_at))),
    CONSTRAINT users_email_valid CHECK (((email)::text ~ '^[^@]+@[^@]+$'::text)),
    CONSTRAINT users_pass_has CHECK ((length((password)::text) = ANY (ARRAY[59, 60]))),
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
-- Name: storage id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.storage ALTER COLUMN id SET DEFAULT nextval('public.storage_id_seq'::regclass);


--
-- Name: subscriptions id; Type: DEFAULT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.subscriptions ALTER COLUMN id SET DEFAULT nextval('public.subscriptions_id_seq'::regclass);


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
1	0	t	2022-04-13 13:38:42.292894
2	20220409194136	t	2022-04-13 13:38:42.300046
\.


--
-- Data for Name: storage; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.storage (id, cid, user_id, name, dir, public, storage_bucket, storage_key, hash, upload_stage, pinned, versions, metadata, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: subscriptions; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.subscriptions (id, user_id, promotion, created_at, updated_at, price, currency, valid_from, valid_to, config) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: animus
--

COPY public.users (id, username, firstname, lastname, email, password, created_at, updated_at, deleted_at) FROM stdin;
1	admin	Animus	Administrator	admin@example.com	$2a$12$IRWQnDUmZ.OjvTuLKaBNte09IcwrvtcPni1G4rBYVpZLW0WJOuPnC	2022-04-13 16:01:20.550377	2022-04-13 16:01:20.551072	\N
\.


--
-- Name: gateways_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.gateways_id_seq', 1, false);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 2, true);


--
-- Name: storage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.storage_id_seq', 1, false);


--
-- Name: subscriptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: animus
--

SELECT pg_catalog.setval('public.subscriptions_id_seq', 1, false);


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
-- Name: gateways gateways_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.gateways
    ADD CONSTRAINT gateways_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: storage storage_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.storage
    ADD CONSTRAINT storage_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: subscriptions subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: animus
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

