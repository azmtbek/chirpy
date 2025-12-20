-- +goose Up
ALTER TABLE chirps
ALTER COLUMN user_id SET NOT NULL;


-- +goose Down
ALTER TABEL chirps
ALTER COLUMN user_id DROP NOT NULL;
