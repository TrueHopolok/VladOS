SELECT id, user_id, data
FROM suggestion
WHERE type = $1
ORDER BY random()
LIMIT 1;