SELECT text 
FROM m8b
WHERE type = ?
ORDER BY random()
LIMIT 1;