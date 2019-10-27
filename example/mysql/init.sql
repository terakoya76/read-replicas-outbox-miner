CREATE DATABASE IF NOT EXISTS read_replicas_outbox_miner_db DEFAULT CHAR SET utf8mb4 DEFAULT collate utf8mb4_general_ci;
CREATE TABLE IF NOT EXISTS read_replicas_outbox_miner_db.progresses (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    database_name VARCHAR(255) NOT NULL DEFAULT "",
    table_name VARCHAR(255) NOT NULL DEFAULT "",
    track_key VARCHAR(255) NOT NULL DEFAULT "",
    position BIGINT NOT NULL DEFAULT 0,
    UNIQUE (database_name, table_name)
);

CREATE DATABASE IF NOT EXISTS outbox_db DEFAULT CHAR SET utf8mb4 DEFAULT collate utf8mb4_general_ci;
CREATE TABLE IF NOT EXISTS outbox_db.outbox_a (
    id INT NOT NULL auto_increment PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    data VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS outbox_db.outbox_b (
    b_id INT NOT NULL auto_increment PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    data VARCHAR(255)
);
