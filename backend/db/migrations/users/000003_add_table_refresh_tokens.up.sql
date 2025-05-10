create schema if not exists users;

create table if not exists users.refresh_tokens(
    id SERIAL PRIMARY KEY,
    user_id INT,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
);