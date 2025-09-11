SELECT u.id, u.firstname, u.username
FROM user AS u
JOIN login AS l
ON u.id = l.user_id
WHERE l.code = $1; 