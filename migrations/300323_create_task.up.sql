CREATE TABLE IF NOT EXISTS task (
    taskId VARCHAR(64) PRIMARY KEY,
    name VARCHAR(512),
    description TEXT,
    status INTEGER,
    createdAt INTEGER NOT NULL,
    updatedAt INTEGER,
    isDeleted INTEGER
);