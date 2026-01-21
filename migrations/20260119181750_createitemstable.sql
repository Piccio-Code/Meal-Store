-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    current_capacity INT NOT NULL,
    store_id INT NOT NULL,
    version uuid NOT NULL DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP,

    CONSTRAINT store_id_fk
        FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
