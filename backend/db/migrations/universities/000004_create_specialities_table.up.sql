CREATE SCHEMA if not exists universities;

CREATE TABLE if not exists universities.specialities(
    id SERIAL PRIMARY KEY,
    university_id INT REFERENCES universities.universities(id),
    name TEXT NOT NULL,
    budget_points FLOAT,
    contract_points FLOAT,
    cost INT NOT NULL
)