-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goals (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    color VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    deadline TIMESTAMP NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    current NUMERIC(15,2) DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1, 
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT unique_user_goal_name UNIQUE (user_id, name)
);

CREATE INDEX IF NOT EXISTS idx_goals_name_tsvector ON goals USING gin(to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS idx_goals_user_id ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_deleted ON goals(deleted);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_goals_deleted;
DROP INDEX IF EXISTS idx_goals_user_id;
DROP INDEX IF EXISTS idx_goals_name_tsvector;
DROP TABLE IF EXISTS goals;
-- +goose StatementEnd
