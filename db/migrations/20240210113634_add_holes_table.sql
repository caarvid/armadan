-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS holes (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  nr integer NOT NULL,
  par integer NOT NULL,
  index integer NOT NULL,
  course_id uuid NOT NULL,
  CONSTRAINT fk_course_id
    FOREIGN KEY (course_id)
      REFERENCES courses(id)
      ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE holes;
-- +goose StatementEnd
