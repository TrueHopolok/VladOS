CREATE TABLE stats_bjack (
    user_id INTEGER PRIMARY KEY,
    games_total INTEGER NOT NULL DEFAULT 0,
    score_current INTEGER NOT NULL DEFAULT 0,
    score_best INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (user_id)
    REFERENCES user(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);