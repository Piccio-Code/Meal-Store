-- +goose Up
-- +goose StatementBegin
ALTER TABLE stores
    ADD COLUMN modified_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stores
    DROP COLUMN modified_at;
-- +goose StatementEnd
