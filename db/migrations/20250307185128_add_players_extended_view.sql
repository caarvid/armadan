-- +goose Up
-- +goose StatementBegin
CREATE VIEW IF NOT EXISTS players_extended AS
SELECT
	p.*,
	u.email,
	cast(coalesce(sum(w.points), 0) as integer) AS points,
	cast(coalesce((
		SELECT h.new_hcp FROM hcp_changes h
		WHERE h.player_id = p.id
		ORDER BY datetime(h.valid_from) DESC
		LIMIT 1
	), 36.0) as real) AS hcp
FROM players p
JOIN users u ON u.id = p.user_id
LEFT JOIN winners w ON w.player_id = p.id
GROUP BY p.id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS players_extended;
-- +goose StatementEnd
