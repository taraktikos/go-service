--Version: 1
--Description: Initial migration
CREATE TABLE IF NOT EXISTS users
(
    id            UUID,
    first_name    text,
    last_name     text,
    phone         text
);
