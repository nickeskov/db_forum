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

CREATE UNLOGGED TABLE IF NOT EXISTS forums
(
    slug           CITEXT PRIMARY KEY UNIQUE NOT NULL,
    title          TEXT                      NOT NULL,
    threads        INTEGER                   NOT NULL,
    posts          BIGINT                    NOT NULL,
    owner_nickname CITEXT                    NOT NULL,

    FOREIGN KEY (owner_nickname) REFERENCES users (nickname)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums_users_nicknames
(
    forum_slug    CITEXT NOT NULL,
    user_nickname CITEXT NOT NULL,

    FOREIGN KEY (forum_slug) REFERENCES forums (slug)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_nickname) REFERENCES users (nickname)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    UNIQUE (forum_slug, user_nickname),
    CONSTRAINT forums_users_nicknames_pk PRIMARY KEY (forum_slug, user_nickname)
);

CREATE VIEW forums_users AS
SELECT fu_nicknames.forum_slug, u.nickname, u.email, u.fullname, u.about
FROM forums_users_nicknames AS fu_nicknames
         JOIN users AS u on fu_nicknames.user_nickname = u.nickname;
