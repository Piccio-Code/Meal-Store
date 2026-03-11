-- +goose Up
-- +goose StatementBegin
CREATE TABLE eatenItems (
    id SERIAL PRIMARY KEY,
    quantity INT NOT NULL,
    eaten_date TIMESTAMP NOT NULL DEFAULT NOW(),
    item_id INT NOT NULL,

    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS eatenItems;
-- +goose StatementEnd
