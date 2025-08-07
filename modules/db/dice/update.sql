INSERT INTO dice (user_id, user_name, throws_total, throws_won, streak_best, streak_current)
VALUES ($1, $2, 1, CASE WHEN $3 = 6 THEN 1 ELSE 0 END, CASE WHEN $3 = 6 THEN 1 ELSE 0 END, CASE WHEN $3 = 6 THEN 1 ELSE 0 END)
ON CONFLICT (user_id) DO UPDATE
SET
    user_name = $2,
    throws_total = dice.throws_total + 1,
    throws_won = dice.throws_won + CASE WHEN $3 = 6 THEN 1 ELSE 0 END,
    streak_current = CASE WHEN $3 = 6 THEN dice.streak_current + 1 ELSE 0 END,
    streak_best = CASE WHEN streak_best <= streak_current AND $3 = 6  THEN dice.streak_current + 1 ELSE streak_best END