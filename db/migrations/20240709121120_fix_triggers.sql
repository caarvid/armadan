-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_player_hcp() RETURNS trigger AS $player_hcp$
  BEGIN
    IF (TG_OP = 'DELETE') THEN
      UPDATE players SET hcp = OLD.old_hcp WHERE id = OLD.player_id;
    ELSIF (TG_OP = 'UPDATE' OR TG_OP = 'INSERT') THEN
      UPDATE players SET hcp = NEW.new_hcp WHERE id = NEW.player_id;
    END IF;

    RETURN NULL;
  END;
$player_hcp$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER player_hcp 
AFTER INSERT OR UPDATE OR DELETE ON rounds
  FOR EACH ROW EXECUTE FUNCTION update_player_hcp();

CREATE OR REPLACE FUNCTION update_player_points() RETURNS trigger AS $player_points$
  BEGIN
    IF (TG_OP = 'DELETE') THEN
      UPDATE players SET points = points - OLD.points WHERE id = OLD.player_id;
    ELSIF (TG_OP = 'INSERT') THEN 
      UPDATE players SET points = points + NEW.points WHERE id = NEW.player_id;
    END IF;

    RETURN NULL;
  END;
$player_points$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER player_points 
AFTER INSERT OR DELETE ON winners
  FOR EACH ROW EXECUTE FUNCTION update_player_points();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER player_hcp;
DROP TRIGGER player_points;
-- +goose StatementEnd
