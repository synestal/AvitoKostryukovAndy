package post

import (
	_ "awesomeProject/internal/app/sqlDAO/get"
	"awesomeProject/internal/app/sqlDAO/models"
	"database/sql"
	"strings"
)

func SetadminState(db *sql.DB, id, state string) error {
	query := `
    INSERT INTO user_tokens(id, token_state)
    VALUES ($1, $2)
    ON CONFLICT (id)
    DO UPDATE SET token_state = EXCLUDED.token_state;
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, state)
	if err != nil {
		return err
	}

	return nil
}

func ChangeHistoryBannersStorage(db *sql.DB, number, id string) (string, error) {
	query := `WITH selected_row AS (
  SELECT *
  FROM history_banenrs
  WHERE id_banner = $2
  LIMIT 1 OFFSET $1-1
)
DELETE FROM history_banenrs AS cte
WHERE id = (SELECT id FROM selected_row)
RETURNING cte.id_banner, cte.tag_list, cte.features_id, cte.title_banner, cte.text_banner, cte.url_banner, cte.banner_state, cte.created_at, cte.updated_at;
	`
	rows, err := db.Query(query, number, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "NULL", nil
	}

	var bannerstable models.FilteredBanner
	tags := make([]uint8, 0, 1)

	err = rows.Scan(
		&bannerstable.Id, &tags, &bannerstable.FeatureIds,
		&bannerstable.Banner.Title, &bannerstable.Banner.Text, &bannerstable.Banner.Url,
		&bannerstable.Flag, &bannerstable.CreatedAt, &bannerstable.UpdatedAt,
	)
	if err != nil {
		return "", err
	}
	str := string(tags)
	str = strings.ReplaceAll(str, "{", "")
	str = strings.ReplaceAll(str, "}", "")
	values := strings.Split(str, ",")
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = strings.TrimSpace(v)
	}
	content := []string{bannerstable.Banner.Title, bannerstable.Banner.Text, bannerstable.Banner.Url}

	ans, err := UpdateBannersStorage(db, result, bannerstable.FeatureIds, bannerstable.Id, bannerstable.Flag, content)
	if err != nil {
		return "", err
	}
	return ans, nil
}
