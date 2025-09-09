CREATE TABLE user_data (
    id INTEGER PRIMARY KEY, /* FK-able with excpetion of conversation table, since it is often cleared and is independent from stats and authefication */
    firstname TEXT NOT NULL,
    username TEXT NOT NULL DEFAULT ""
);