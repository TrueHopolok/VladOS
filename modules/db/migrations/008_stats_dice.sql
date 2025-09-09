CREATE TABLE stats_dice (
    user_id INTEGER PRIMARY KEY,
    games_total INTEGER NOT NULL DEFAULT 0,
    score_current INTEGER NOT NULL DEFAULT 0,
    score_best INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (user_id)
    REFERENCES user_data(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);