CREATE TABLE IF NOT EXISTS tasks (
    id          text PRIMARY KEY,
    title       text,
    description text,
    status      text        NOT NULL DEFAULT 'TaskStatus_TODO',
    user_id     text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    deleted_at  timestamptz,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
