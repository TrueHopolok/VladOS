CREATE TABLE suggestion (
    user_id INTEGER,
    type TEXT NOT NULL,
    suggestion BLOB,
    
    FOREIGN KEY (user_id)
    REFERENCES user(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
);