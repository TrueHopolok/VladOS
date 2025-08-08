SELECT text 
FROM m8b
WHERE type = $1
ORDER BY random()
LIMIT 1;