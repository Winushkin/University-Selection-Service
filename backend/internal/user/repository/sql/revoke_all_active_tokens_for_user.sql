DELETE
FROM users.refresh_tokens
WHERE user_id = $1;