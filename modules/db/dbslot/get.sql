SELECT throws_total, score_current, score_best
FROM dice
WHERE user_id = $1;