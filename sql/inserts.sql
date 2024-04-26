CREATE TABLE doings(
    id bigserial primary key ,
    chat_id bigint,
    name text,
    importance integer,
    time timestamp,
    status boolean
);
CREATE TABLE languages(
    chat_id bigint,
    language text
);


DROP TABLE doings, languages;


-- CREATE TABLE doings(
--                        id bigserial primary key ,
--                        chat_id bigint,
--                        name text,
--                        importance integer,
--                        time timestamp,
--                        status boolean default false
-- );
-- CREATE TABLE languages(
--                           chat_id bigint,
--                           language text
-- );