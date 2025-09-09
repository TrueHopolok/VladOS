CREATE TABLE conversation (
    user_id INTEGER PRIMARY KEY, /* Not FK of user, since table is often cleared and is independent from stats and authefication */
    available INTEGER CHECK(available IN (0, 1)) NOT NULL DEFAULT 1,
	command_name TEXT DEFAULT "",
    additional_data BLOB
);