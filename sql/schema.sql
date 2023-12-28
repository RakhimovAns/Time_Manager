ALTER TABLE doings
DROP COLUMN status,
ADD  COLUMN status boolean  default false;