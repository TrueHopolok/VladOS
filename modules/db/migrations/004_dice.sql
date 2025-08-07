CREATE TABLE dice (
    user_id INTEGER PRIMARY KEY,
    user_name TEXT NOT NULL,
    throws_total INTEGER NOT NULL DEFAULT 0,
    throws_won INTEGER NOT NULL DEFAULT 0,
    streak_best INTEGER NOT NULL DEFAULT 0,
    streak_current INTEGER NOT NULL DEFAULT 0
);