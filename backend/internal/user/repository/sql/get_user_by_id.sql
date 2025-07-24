SELECT id,
       login,
       password,
       coalesce(ege, 0),
       coalesce(speciality, ''),
       coalesce(region, ''),
       coalesce(financing, '')
FROM users.users
WHERE id = $1;