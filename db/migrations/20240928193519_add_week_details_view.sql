-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW week_details AS
SELECT
  w.id,
  w.nr,
  w.finals_date,
  w.is_finals,
  c.id as course_id,
  c.name as course_name,
  t.id as tee_id,
  t.name as tee_name
FROM weeks w
LEFT JOIN courses c ON c.id = w.course_id
LEFT JOIN tees t ON t.id = w.tee_id
GROUP BY w.id, c.id, c.name, t.id, t.name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS week_details;
-- +goose StatementEnd
