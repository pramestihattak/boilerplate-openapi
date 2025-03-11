-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD COLUMN phone_number VARCHAR(15);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP COLUMN phone_number;
