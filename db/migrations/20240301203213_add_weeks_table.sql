-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS weeks (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  nr integer NOT NULL UNIQUE,
  course_id uuid NOT NULL,
  tee_id uuid NOT NULL,
  CONSTRAINT fk_course_id
    FOREIGN KEY (course_id)
    REFERENCES courses(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_tee_id
    FOREIGN KEY (tee_id)
    REFERENCES tees(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE weeks;
-- +goose StatementEnd
