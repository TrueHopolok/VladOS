WITH ranked AS (
    SELECT
        s.user_id,
        u.firstname,
        u.username,
        s.games_total,
        s.score_current,
        s.score_best,
        RANK() OVER (ORDER BY s.score_best DESC) AS rank,
        COUNT(*) OVER () AS players_total
    FROM %s as s
    JOIN user AS u
    ON s.user_id = u.id
)
SELECT * FROM ranked
WHERE rank <= 10;