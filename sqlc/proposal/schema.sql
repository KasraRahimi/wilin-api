CREATE TABLE IF NOT EXISTS proposals (
    id int PRIMARY KEY AUTO_INCREMENT,
    user_id int,
    entry varchar(255) NOT NULL,
    pos varchar(255) NOT NULL,
    gloss varchar(255) NOT NULL,
    notes varchar(2047) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);