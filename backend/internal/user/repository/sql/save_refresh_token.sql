INSERT INTO users.refresh_tokens
    (user_id, token, expires_at)
VALUES ($1, $2, $3);