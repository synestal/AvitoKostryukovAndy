package post

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type BannerId struct {
	ID string `json:"banner_id"`
}

func CreateNemBannerStorage(db *sql.DB, tags []string, feature, active string, content []string) (*BannerId, error) {
	query := `
		INSERT INTO banners_storage(tag_list, features_id, title_banner, text_banner, url_banner, banner_state)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_banner
	`

	var bannerId BannerId
	err := db.QueryRow(query, pq.Array(tags), feature, content[0], content[1], content[2], active).Scan(&bannerId.ID)
	if err != nil {
		return nil, err
	}

	return &bannerId, nil
}

func UpdateBannersStorage(db *sql.DB, tags []string, feature, id, active string, content []string) (string, error) {
	query := `WITH updated_rows AS (
    UPDATE banners_storage
    SET tag_list = $6,
        features_id = $7,
        banner_state = $2,
        title_banner = $3,
        text_banner = $4,
        url_banner = $5
    WHERE id_banner = $1
    RETURNING *
)
SELECT CASE 
           WHEN EXISTS (SELECT 1 FROM updated_rows) THEN 'updated' 
           ELSE 'no banner' 
       END AS result;
    `

	rows, err := db.Query(query, id, active, content[0], content[1], content[2], pq.Array(tags), feature)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var ans string
	if rows.Next() {
		err := rows.Scan(&ans)
		if err != nil {
			return "", err
		}
	}
	fmt.Println(ans)

	return ans, nil
}

func DeleterBanners(db *sql.DB, item string) (string, error) {
	query := `
		DELETE FROM banners_storage 
		WHERE id_banner = $1
		RETURNING 1 AS deleted
	`

	var deleted int
	err := db.QueryRow(query, item).Scan(&deleted)
	switch {
	case err == sql.ErrNoRows:
		return "not found", nil
	case err != nil:
		return "", err
	}

	if deleted > 0 {
		return "deleted", nil
	} else {
		return "not found", nil
	}
}

func DeleterBannersPostponed(db *sql.DB, item []int) error {
	query := `
		INSERT INTO delayed_deletions (id_item)
		SELECT unnest($1::int[])
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pq.Array(item))
	if err != nil {
		return err
	}

	return nil
}
