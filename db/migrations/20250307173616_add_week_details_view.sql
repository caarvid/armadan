-- +goose Up
-- +goose StatementBegin
CREATE VIEW IF NOT EXISTS week_details AS
SELECT
  w.*,
  c.name as course_name,
  t.name as tee_name
FROM weeks w
JOIN courses c ON c.id = w.course_id
JOIN tees t ON t.id = w.tee_id
GROUP BY w.id, c.id, c.name, t.id, t.name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS week_details;
-- +goose StatementEnd

