BEGIN TRANSACTION;

CREATE TABLE tokens(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token TEXT UNIQUE NOT NULL
);

UPDATE authors SET author = 'Signe Soolalõhe' WHERE author = 'Signe Soolõhe';

COMMIT;
