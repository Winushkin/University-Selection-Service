CREATE TABLE if not exists universities.budget(
    id SERIAL PRIMARY KEY,
    prestige INT NOT NULL,
    name TEXT NOT NULL,
    rank FLOAT NOT NULL,
    speciality TEXT NOT NULL,
    quality SMALLINT NOT NULL,
    points FLOAT NOT NULL,
    dormitory BOOLEAN,
    labs BOOLEAN,
    sport BOOLEAN,
    scholarship INT,
    region_name TEXT NOT NULL
);


