CREATE SCHEMA if not exists universities;

CREATE TABLE if not exists universities.regions(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);