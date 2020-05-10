DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    nickname CITEXT PRIMARY KEY UNIQUE NOT NULL,
    email    CITEXT UNIQUE             NOT NULL,
    fullname TEXT                      NOT NULL,
    about    TEXT                      NOT NULL
);
