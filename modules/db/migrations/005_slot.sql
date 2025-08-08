CREATE TABLE slot (
    user_id INTEGER PRIMARY KEY,
    throws_total INTEGER NOT NULL DEFAULT 0,
    score_current INTEGER NOT NULL DEFAULT 0,
    score_best INTEGER NOT NULL DEFAULT 0
);