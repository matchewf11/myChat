CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE CHECK(length(username) <= 16),
    password TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    last_login TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_name TEXT NOT NULL UNIQUE CHECK(length(room_name) <= 16),
    room_password TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);

-- CREATE TABLE IF NOT EXISTS user_rooms (
--     user_id INTEGER NOT NULL,
--     room_id INTEGER NOT NULL,
--     PRIMARY KEY (user_id, room_id),
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
--     FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
-- );

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    body TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    author INTEGER NOT NULL,
    FOREIGN KEY (author) REFERENCES users(id) ON DELETE CASCADE
-- room INTEGER NOT NULL,
-- FOREIGN KEY (room) REFERENCES rooms(id) ON DELETE CASCADE
);

-- CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author);
-- CREATE INDEX IF NOT EXISTS idx_posts_room ON posts(room);
-- CREATE INDEX IF NOT EXISTS idx_user_rooms_user ON user_rooms(user_id);
-- CREATE INDEX IF NOT EXISTS idx_user_rooms_room ON user_rooms(room_id);
