CREATE TABLE doings(
    id bigserial primary key ,
    chat_id integer not null,
    name text not null ,
    importance integer not null ,
    time timestamp not null
)