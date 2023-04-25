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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ip_addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ip_addresses (
    id integer NOT NULL,
    subnet_id integer NOT NULL,
    ip_address bigint NOT NULL,
    hostname character varying(50) NOT NULL
);


--
-- Name: ip_addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ip_addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ip_addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ip_addresses_id_seq OWNED BY public.ip_addresses.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: subnets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subnets (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    network_address bigint NOT NULL,
    mask_length integer NOT NULL,
    gateway bigint,
    name_server bigint,
    description text
);


--
-- Name: subnets_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.subnets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: subnets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.subnets_id_seq OWNED BY public.subnets.id;


--
-- Name: ip_addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ip_addresses ALTER COLUMN id SET DEFAULT nextval('public.ip_addresses_id_seq'::regclass);


--
-- Name: subnets id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subnets ALTER COLUMN id SET DEFAULT nextval('public.subnets_id_seq'::regclass);


--
-- Name: ip_addresses ip_addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ip_addresses
    ADD CONSTRAINT ip_addresses_pkey PRIMARY KEY (id);


--
-- Name: ip_addresses ip_addresses_subnet_id_ip_address_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ip_addresses
    ADD CONSTRAINT ip_addresses_subnet_id_ip_address_key UNIQUE (subnet_id, ip_address);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: subnets subnets_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subnets
    ADD CONSTRAINT subnets_name_key UNIQUE (name);


--
-- Name: subnets subnets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subnets
    ADD CONSTRAINT subnets_pkey PRIMARY KEY (id);


--
-- Name: ip_addresses ip_addresses_subnet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ip_addresses
    ADD CONSTRAINT ip_addresses_subnet_id_fkey FOREIGN KEY (subnet_id) REFERENCES public.subnets(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20230402112853'),
    ('20230410010441');
