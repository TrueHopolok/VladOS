SELECT games_total, score_current, score_best
FROM %s
WHERE user_id = $1;