CREATE TABLE IF NOT EXISTS progresses (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    database_name VARCHAR(255) NOT NULL DEFAULT "",
    table_name VARCHAR(255) NOT NULL DEFAULT "",
    track_key VARCHAR(255) NOT NULL DEFAULT "",
    position BIGINT NOT NULL DEFAULT 0,
    UNIQUE (database_name, table_name)
);