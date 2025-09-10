SELECT u.id, u.firstname, u.username
FROM user_data AS u
JOIN user_login AS l
ON u.id = l.user_id
WHERE l.code = $1; 