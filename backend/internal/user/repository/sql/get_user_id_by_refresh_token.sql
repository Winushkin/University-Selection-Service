SELECT user_id
FROM users.refresh_tokens
WHERE token = $1;