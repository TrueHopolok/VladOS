SELECT throws_total, score_current, score_best
FROM slot
WHERE user_id = $1;