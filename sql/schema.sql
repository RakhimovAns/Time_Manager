ALTER TABLE doings
ADD  COLUMN status boolean not null default false,
ADD COLUMN done_time timestamp  ;