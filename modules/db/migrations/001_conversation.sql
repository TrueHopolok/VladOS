CREATE TABLE conversation (
    user_id INTEGER PRIMARY KEY,
    available INTEGER CHECK(available IN (0, 1)) NOT NULL DEFAULT 1,
	command_name TEXT DEFAULT "",
    additional_data BLOB
);