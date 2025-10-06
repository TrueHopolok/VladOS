SELECT pun 
FROM pun
WHERE suffix=$1
ORDER BY random()
LIMIT 1;