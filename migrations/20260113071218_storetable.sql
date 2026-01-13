-- +goose Up
-- +goose StatementBegin
CREATE TABLE store (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    UserID VARCHAR(255),
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES "user"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS store;
-- +goose StatementEnd
