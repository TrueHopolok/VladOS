CREATE TABLE login (
    user_id INTEGER,
    code TEXT,
    expiration INTEGER,

    PRIMARY KEY (user_id, code),

    FOREIGN KEY (user_id)
    REFERENCES user(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);