-- +goose Up
-- +goose StatementBegin
INSERT INTO "user"(id, name, email, "emailVerified", image)
VALUES ('test', 'testUser', 'test@test.com', false, 'test');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "user"
WHERE id='test'
-- +goose StatementEnd
