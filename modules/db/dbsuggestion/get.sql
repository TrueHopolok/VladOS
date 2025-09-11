SELECT id, user_id, type, data
FROM suggestion
ORDER BY random()
LIMIT 1;