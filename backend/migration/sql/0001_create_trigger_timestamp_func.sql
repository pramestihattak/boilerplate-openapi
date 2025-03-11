-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp_utc()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated = NOW() at time zone 'utc';
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION trigger_set_timestamp_utc();
