-- +goose Up
-- +goose StatementBegin
ALTER TABLE stores
    ADD CONSTRAINT stores_user_id_name_unique
        UNIQUE (user_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stores
    DROP CONSTRAINT stores_user_id_name_unique;
-- +goose StatementEnd
