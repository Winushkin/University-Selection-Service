create schema if not exists schema_name;

create table if not exists schema_name.users (
    Id SERIAL PRIMARY KEY,
    Login TEXT UNIQUE NOT NULL,
    Password TEXT,
    Ege INT,
    Gpa FLOAT,
    Speciality TEXT,
    EduType TEXT,
    Town TEXT,
    Financing TEXT
);

create table if not exists schema_name.refresh_tokens
(
    id SERIAL PRIMARY KEY,
    user_id INT,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
);