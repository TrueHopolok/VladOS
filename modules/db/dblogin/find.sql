SELECT u.id, u.firstname, u.username, EXISTS(
    SELECT *
    FROM admin AS a
    WHERE a.user_id = u.id
) AS admin
FROM user AS u
JOIN login AS l
ON u.id = l.user_id
WHERE l.code = $1; 