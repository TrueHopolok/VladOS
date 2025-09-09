CREATE TABLE user_login (
    user_id INTEGER PRIMARY KEY,
    code TEXT,
    expiration INTEGER,

    FOREIGN KEY (user_id)
    REFERENCES user(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);