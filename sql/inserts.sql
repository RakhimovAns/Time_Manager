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
