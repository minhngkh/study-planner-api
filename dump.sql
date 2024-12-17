PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE,
    password TEXT,
    google_id TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO user VALUES(1,'minhngkh@gmail.com','$2a$10$3dQ0jbHKdlzYDzWm7zUiUugCbBs28OuK6A4hMtGyZKYFKvSOAOpae',NULL,'2024-12-17 19:44:56.575185183+07:00','2024-12-17 19:44:56.575185183+07:00');
CREATE TABLE IF NOT EXISTS user_session (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP,
  
    FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE
);
INSERT INTO user_session VALUES(1,1,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNzM1MDQ0Mjk2fQ.n3OwQYOlYshSBtJsfFxY6HDOjdRRwq6PtM9I_ClE0WQ','2024-12-24 19:44:56.747551461+07:00','2024-12-17 19:44:56.829644885+07:00');
CREATE TABLE IF NOT EXISTS task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    name TEXT NOT NULL,
    description TEXT,
    priority TEXT NOT NULL, -- Allowed values: 'Low', 'Medium', 'High' (validate in application)
    estimated_time INTEGER, -- in minutes
    status TEXT NOT NULL,   -- Allowed values: 'Todo', 'In Progress', 'Completed', 'Expired' (validate in application)
    start_time TIMESTAMP,   -- When the task is scheduled to start
    end_time TIMESTAMP,     -- When the task is scheduled to end
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE
);
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('user',1);
INSERT INTO sqlite_sequence VALUES('user_session',1);
COMMIT;
