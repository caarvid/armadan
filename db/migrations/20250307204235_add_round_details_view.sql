-- +goose Up
-- +goose StatementBegin
CREATE VIEW IF NOT EXISTS full_rounds AS
SELECT
	r.*,
	rd.net_in,
	rd.net_out,
	rd.net_total,
	rd.gross_in,
	rd.gross_out,
	rd.gross_total,
	hc.old_hcp,
	hc.new_hcp,
	coalesce((
		SELECT json_group_array(
			json_object(
				'id', s.id,
				'strokes', s.strokes,
				'round_id', s.round_id,
				'hole', json_object(
					'id', h.id,
					'nr', h.nr,
					'par', h.par,
					'stroke_index', h.stroke_index
				)
			)
		) 
		FROM scores s 		
		JOIN holes h ON h.id = s.hole_id
		WHERE s.round_id = r.id
		ORDER BY h.nr ASC
	), '[]') AS scores
FROM rounds r
JOIN round_details rd ON rd.round_id = r.id
JOIN hcp_changes hc ON hc.round_id = r.id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS full_rounds;
-- +goose StatementEnd

