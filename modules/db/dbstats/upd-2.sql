INSERT INTO %s (user_id, games_total, score_current, score_best)
VALUES (
    $1, 
    1, 
    CASE WHEN $2 > 0 THEN $2 ELSE 0 END,
    CASE WHEN $2 > 0 THEN $2 ELSE 0 END
)
ON CONFLICT (user_id) DO UPDATE
SET
    games_total = games_total + 1,
    score_current = CASE WHEN $2 > 0 THEN score_current + $2 ELSE 0 END,
    score_best = CASE WHEN $2 > 0 AND score_current + $2 > score_best THEN score_current + $2 ELSE score_best END