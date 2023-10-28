CREATE ROLE app WITH
	LOGIN
	NOSUPERUSER
    CREATEDB
	CREATEROLE
	NOINHERIT
	NOREPLICATION
	CONNECTION LIMIT -1
	PASSWORD 'APP';

CREATE DATABASE chat
    WITH
    OWNER = app
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

\c chat

-- DROP TABLE IF EXISTS public.transaction;
-- CREATE TABLE public.transaction
-- (
--     id serial NOT NULL,
--     item_desc text,
--     purchase_date date,
--     transaction_id character varying(36),
--     amount double precision,
--     PRIMARY KEY (id)
-- );

-- ALTER TABLE IF EXISTS public.transaction
--     OWNER to app;