-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS results (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  has_results boolean NOT NULL DEFAULT false,
  week_id uuid NOT NULL,
  CONSTRAINT fk_week_id
    FOREIGN KEY (week_id)
      REFERENCES weeks(id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rounds (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  net_in integer NOT NULL,
  net_out integer NOT NULL,
  gross_in integer NOT NULL,
  gross_out integer NOT NULL,
  old_hcp decimal(3, 1) NOT NULL,
  new_hcp decimal(3, 1) NOT NULL,
  player_id uuid NOT NULL,
  result_id uuid NOT NULL,
  CONSTRAINT fk_player_id
    FOREIGN KEY (player_id)
      REFERENCES players(id)
      ON DELETE CASCADE,
  CONSTRAINT fk_result_id
    FOREIGN KEY (result_id)
      REFERENCES results(id)
      ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_player_hcp() RETURNS trigger AS $player_hcp$
  BEGIN
    IF (TG_OP = 'DELETE') THEN
      UPDATE players p SET p.hcp = OLD.old_hcp WHERE p.id = OLD.player_id;
    ELSIF (TG_OP = 'UPDATE' OR TG_OP = 'INSERT') THEN
      UPDATE players p SET p.hcp = NEW.new_hcp WHERE p.id = NEW.player_id;
    END IF;

    RETURN NULL;
  END;
$player_hcp$ LANGUAGE plpgsql;

CREATE TRIGGER player_hcp 
AFTER INSERT OR UPDATE OR DELETE ON rounds
  FOR EACH ROW EXECUTE FUNCTION update_player_hcp();

CREATE TABLE IF NOT EXISTS winners (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  points integer NOT NULL,
  result_id uuid NOT NULL,
  player_id uuid NOT NULL,
  CONSTRAINT fk_player_id
    FOREIGN KEY (player_id)
      REFERENCES players(id)
      ON DELETE CASCADE,
  CONSTRAINT fk_result_id
    FOREIGN KEY (result_id)
      REFERENCES results(id)
      ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_player_points() RETURNS trigger AS $player_points$
  BEGIN
    IF (TG_OP = 'DELETE') THEN
      UPDATE players p SET p.points = p.points - OLD.points WHERE p.id = OLD.player_id;
    ELSIF (TG_OP = 'INSERT') THEN 
      UPDATE players p SET p.points = p.points + NEW.points WHERE p.id = NEW.player_id;
    END IF;

    RETURN NULL;
  END;
$player_points$ LANGUAGE plpgsql;

CREATE TRIGGER player_points 
AFTER INSERT OR DELETE ON winners
  FOR EACH ROW EXECUTE FUNCTION update_player_points();

CREATE TABLE IF NOT EXISTS scores (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  strokes integer NOT NULL,
  hole_id uuid NOT NULL,
  round_id uuid NOT NULL,
  CONSTRAINT fk_hole_id
    FOREIGN KEY (hole_id)
      REFERENCES holes(id)
      ON DELETE CASCADE,
  CONSTRAINT fk_round_id
    FOREIGN KEY (round_id)
      REFERENCES rounds(id)
      ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE results;
DROP TABLE rounds;
DROP TABLE winners;
DROP TABLE scores;
-- +goose StatementEnd
