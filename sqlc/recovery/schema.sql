CREATE TABLE IF NOT EXISTS recoveries (
    id varchar(255) PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);