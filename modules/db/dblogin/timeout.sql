DELETE FROM login
WHERE (
    $1
    AND
    user_id = $2
)
OR expiration < $3;