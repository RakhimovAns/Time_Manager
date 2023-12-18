CREATE TABLE doings(
                       id bigserial primary key ,
                       chat_id bigint not null,
                       name text not null ,
                       importance integer not null ,
                       time timestamp not null
)