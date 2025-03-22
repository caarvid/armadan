-- +goose Up
-- +goose StatementBegin
CREATE VIEW IF NOT EXISTS course_details AS
SELECT
	c.*,
	coalesce((
		SELECT json_group_array(
			json_object(
				'id', h.id,
				'nr', h.nr,
				'par', h.par,
				'index', h.stroke_index,
				'course_id', h.course_id
			)
		) FROM holes h WHERE h.course_id = c.id
	), '[]') AS holes,
	coalesce((
		SELECT json_group_array(
			json_object(
				'id', t.id,
				'name', t.name,
				'cr', t.cr,
				'slope', t.slope,
				'course_id', t.course_id
			)
		) FROM tees t WHERE t.course_id = c.id
	), '[]') AS tees
FROM courses c;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS course_details;
-- +goose StatementEnd
