UPDATE users.users
SET ege = $1,
    speciality = $2,
    region = $3,
    financing = $4
WHERE id = $5;