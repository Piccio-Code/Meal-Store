-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE stores
    ADD COLUMN version uuid NOT NULL DEFAULT uuid_generate_v4();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "uuid-ossp";

ALTER TABLE stores
    DROP COLUMN version;
-- +goose StatementEnd
