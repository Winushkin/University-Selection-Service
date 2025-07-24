INSERT INTO users.users (login, password)
VALUES ($1, $2) RETURNING id