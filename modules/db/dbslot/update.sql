INSERT INTO slot (user_id, throws_total, score_current, score_best)
VALUES (
    $1, 
    1, 
    CASE WHEN $2 > 0 THEN $2 ELSE 0 END,
    CASE WHEN $2 > 0 THEN $2 ELSE 0 END
)
ON CONFLICT (user_id) DO UPDATE
SET
    throws_total = slot.throws_total + 1,
    score_current = CASE WHEN $2 > 0 THEN slot.score_current + $2 ELSE 0 END,
    score_best = CASE WHEN $2 > 0 AND slot.score_current + $2 > slot.score_best THEN slot.score_current + $2 ELSE slot.score_best END