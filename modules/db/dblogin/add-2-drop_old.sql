DELETE FROM user_login
WHERE user_id = $1
OR expiration < $2;