WITH ranked AS (
    SELECT
        user_id,
        games_total,
        score_current,
        score_best,
        RANK() OVER (ORDER BY score_best DESC) AS rank,
        COUNT(*) OVER () AS players_total
    FROM %s
)
SELECT * FROM ranked
WHERE rank <= 3
UNION
SELECT * FROM ranked
WHERE user_id = $1
ORDER BY rank;