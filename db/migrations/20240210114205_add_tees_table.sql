-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tees (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  slope integer NOT NULL,
  cr numeric(3, 1) NOT NULL,
  course_id uuid NOT NULL,
  CONSTRAINT fk_course_id
    FOREIGN KEY (course_id)
      REFERENCES courses(id)
      ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tees;
-- +goose StatementEnd
