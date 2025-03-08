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
	hc.new_hcp
FROM rounds r
JOIN round_details rd ON rd.round_id = r.id
JOIN hcp_changes hc ON hc.round_id = r.id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS full_rounds;
-- +goose StatementEnd

