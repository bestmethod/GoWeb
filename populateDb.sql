CREATE TABLE session (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    session_id string(50),
    session_key string(130),
    expires INTEGER,
    keep_logged_in bool
);

CREATE TABLE users (
    user_id INTEGER PRIMARY KEY,
    username string(50),
    password string(50),
    registered INTEGER,
    last_login INTEGER
);

