-- +goose Up
CREATE TABLE links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_id INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    max_downloads INTEGER,
    download_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    expires_at DATETIME,
    revoked_at DATETIME
);

-- +goose Down
DROP TABLE links;