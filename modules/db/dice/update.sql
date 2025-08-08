INSERT INTO dice (user_id, throws_total, score_current, score_best)
VALUES (
    $1, 
    1, 
    CASE WHEN $2 > 1 THEN $2 ELSE 0 END,
    CASE WHEN $2 > 1 THEN $2 ELSE 0 END
)
ON CONFLICT (user_id) DO UPDATE
SET
    throws_total = dice.throws_total + 1,
    score_current = CASE WHEN $2 > 1 THEN dice.score_current + $2 ELSE 0 END,
    score_best = CASE WHEN $2 > 1 AND dice.score_current + $2 > dice.score_best THEN dice.score_current + $2 ELSE dice.score_best END