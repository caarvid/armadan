-- +goose Up
-- +goose StatementBegin
CREATE VIEW IF NOT EXISTS players_with_points AS
SELECT
	p.*,
	u.email,
	cast(coalesce(sum(w.points), 0) as integer) AS points
FROM players p
JOIN users u ON u.id = p.user_id
LEFT JOIN winners w ON w.player_id = p.id
GROUP BY p.id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS players_with_points;
-- +goose StatementEnd
