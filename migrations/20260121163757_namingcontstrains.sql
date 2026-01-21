-- +goose Up
-- +goose StatementBegin
ALTER TABLE items
    ADD CONSTRAINT name_store_id_unique UNIQUE(store_id, name);

ALTER TABLE stores
    ADD CONSTRAINT name_user_unique UNIQUE(user_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items
    DROP CONSTRAINT name_store_id_unique;

ALTER TABLE stores
    DROP CONSTRAINT name_user_unique;
-- +goose StatementEnd
