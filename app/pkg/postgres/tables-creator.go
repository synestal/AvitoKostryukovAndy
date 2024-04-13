package postgres

import "database/sql"

func CreateTable(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS banners_storage(
    id_banner  SERIAL,
    tag_list INTEGER[],
    features_id INTEGER,
    title_banner TEXT,
    text_banner TEXT,
    url_banner TEXT,
    banner_state BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_banner)
);

CREATE TABLE IF NOT EXISTS delayed_deletions(
    id  SERIAL,
    id_item INTEGER,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user_tokens(
    id INTEGER PRIMARY KEY,
    token_state  BOOLEAN,
    CONSTRAINT unique_id UNIQUE (id)
);

CREATE TABLE IF NOT EXISTS history_banenrs(
    id SERIAL,
    id_banner  SERIAL,
    tag_list INTEGER[],
    features_id INTEGER,
    title_banner TEXT,
    text_banner TEXT,
    url_banner TEXT,
    banner_state BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
		`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func CreateTrigger(db *sql.DB) error {
	query := `
CREATE OR REPLACE FUNCTION process_delayed_deletions()
RETURNS TRIGGER AS $$
BEGIN
    EXECUTE format('DELETE FROM banners_storage WHERE id_banner = ANY(SELECT id_item FROM delayed_deletions)');
    DELETE FROM delayed_deletions;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS before_update_or_insert_on_delayed_deletions ON delayed_deletions;

CREATE TRIGGER before_update_or_insert_on_delayed_deletions
AFTER INSERT ON delayed_deletions
FOR EACH ROW
EXECUTE FUNCTION process_delayed_deletions();
------------------------------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_history_banner()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO history_banenrs (id_banner, tag_list, features_id, title_banner, text_banner, url_banner, banner_state, created_at, updated_at)
    VALUES (NEW.id_banner, NEW.tag_list, NEW.features_id, NEW.title_banner, NEW.text_banner, NEW.url_banner, NEW.banner_state, NEW.created_at, NEW.updated_at);

    DELETE FROM history_banenrs
    WHERE id_banner = NEW.id_banner
    AND id IN (
        SELECT id
        FROM history_banenrs
        WHERE id_banner = NEW.id_banner
        ORDER BY id DESC
        OFFSET 4
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_history_banner_trigger ON banners_storage;

CREATE TRIGGER update_history_banner_trigger
AFTER INSERT OR UPDATE ON banners_storage
FOR EACH ROW
EXECUTE FUNCTION update_history_banner();

------------------------------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_updated_at_trigger ON banners_storage;

CREATE TRIGGER update_updated_at_trigger
BEFORE UPDATE ON banners_storage
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
		`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
