CREATE SCHEMA if not exists universities;
CREATE TABLE if not exists universities.universities(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    site TEXT,
    prestige INT NOT NULL UNIQUE,
    rank FLOAT NOT NULL,
    quality SMALLINT NOT NULL,
    scholarship INT,
    dormitory BOOLEAN,
    labs BOOLEAN,
    sport BOOLEAN,
    region_id INT REFERENCES universities.regions(id)
);


