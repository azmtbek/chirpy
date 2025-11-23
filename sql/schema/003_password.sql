-- +goose Up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS password TEXT DEFAULT 'unset';

ALTER TABLE users
ALTER COLUMN password SET NOT NULL;


-- +goose Down
ALTER TABEL users
DROP COLUMN IF EXISTS password;
