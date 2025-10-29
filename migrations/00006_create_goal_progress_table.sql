-- +goose Up
-- +goose StatementBegin
CREATE TABLE goal_progress (
    id BIGSERIAL PRIMARY KEY,
    goal_id BIGINT NOT NULL REFERENCES goals(id),
    amount NUMERIC(12,2) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    version INT DEFAULT 1,
    deleted BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goal_progress;
-- +goose StatementEnd
