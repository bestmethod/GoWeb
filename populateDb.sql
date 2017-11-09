CREATE TABLE session (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    session_id string(50),
    session_key string(130),
    expires INTEGER,
    keep_logged_in bool
);

