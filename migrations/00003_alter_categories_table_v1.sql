-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories
ADD CONSTRAINT unique_user_category_name UNIQUE (user_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories
DROP CONSTRAINT IF EXISTS unique_user_category_name;
-- +goose StatementEnd
