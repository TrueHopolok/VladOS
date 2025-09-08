SELECT score_best, COUNT(*) FROM %s
GROUP BY score_best
ORDER BY score_best DESC;