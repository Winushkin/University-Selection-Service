create schema if not exists users;

create table if not exists users.users (
     id SERIAL PRIMARY KEY,
     login TEXT UNIQUE NOT NULL,
     password TEXT,
     ege INT,
     speciality TEXT,
     region TEXT,
     financing TEXT
);
