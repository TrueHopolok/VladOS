SELECT throws_total, throws_won, streak_current, streak_best
FROM dice
WHERE user_id = $1;