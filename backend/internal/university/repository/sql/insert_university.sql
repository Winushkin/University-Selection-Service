INSERT INTO universities.universities
    (name,
    site,
    prestige,
    rank,
    quality,
    scholarship,
    dormitory,
    labs,
    sport,
    region_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);