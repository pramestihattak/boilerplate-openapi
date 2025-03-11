-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF not EXISTS users(
  "user_id" UUID DEFAULT uuid_generate_v1() PRIMARY KEY,
  "full_name" VARCHAR(60) NOT NULL,
  "email" VARCHAR(60) NOT NULL,
  "password" VARCHAR(72) NOT NULL,
  "avatar_url" VARCHAR(255) NOT NULL DEFAULT '',
  "verified" BOOLEAN NOT NULL DEFAULT false,
  "verification_token" VARCHAR(10) NOT NULL DEFAULT '',
  "created" timestamp null DEFAULT now(),
  "updated" timestamp null,
  unique(user_id, email)
);

CREATE TRIGGER set_timestamp_utc
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp_utc();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
