package get

import (
	"awesomeProject/internal/app/sqlDAO/models"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type BannerState struct {
	state string
}

func GetBannerFromDB(db *sql.DB, tagID, featureID string) (*models.Banner, string, error) {
	query := `
		SELECT title_banner, text_banner, url_banner, banner_state
		FROM banners_storage
		WHERE features_id = $2 AND $1 = ANY(tag_list);
	`

	rows, err := db.Query(query, tagID, featureID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var banner models.Banner
	var state string
	if rows.Next() {
		err := rows.Scan(&banner.Title, &banner.Text, &banner.Url, &state)
		if err != nil {
			return nil, "", err
		}
	}

	return &banner, state, nil
}

type Admin struct {
	avaliable string
	admin     string
}

func GetAdminState(db *sql.DB, token string) (bool, bool, error) {
	query := `
		SELECT 
			CASE 
				WHEN EXISTS (SELECT 1 FROM user_tokens WHERE id = $1) THEN TRUE
				ELSE FALSE
			END AS found_subj,
			EXISTS (SELECT 1 FROM user_tokens WHERE id = $1 AND token_state = TRUE) AS status
	`

	var foundSubj, status bool
	err := db.QueryRow(query, token).Scan(&foundSubj, &status)
	if err != nil {
		return false, false, err
	}

	return foundSubj, status, nil
}

func GetBannerIdByTag(db *sql.DB, tag, limit, offset string) ([]int, error) {
	var ids []int

	// Подготовка базового запроса
	query := `
		SELECT id_banner
		FROM banners_storage
		WHERE $1 = ANY(tag_list)
	`

	// Добавление условий LIMIT и OFFSET, если они заданы
	if limit != "" {
		query += " LIMIT " + limit
	}
	if offset != "" {
		query += " OFFSET " + offset
	}

	// Подготовка запроса и выполнение
	rows, err := db.Query(query, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Итерация по результатам запроса
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func GetBannerIdByFeature(db *sql.DB, feature, limit, offset string) ([]int, error) {
	var ids []int

	// Подготовка базового запроса
	query := `
		SELECT id_banner
		FROM banners_storage
		WHERE $1 = features_id
	`

	// Добавление условий LIMIT и OFFSET, если они заданы
	if limit != "" {
		query += " LIMIT " + limit
	}
	if offset != "" {
		query += " OFFSET " + offset
	}

	// Подготовка запроса и выполнение
	rows, err := db.Query(query, feature)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Итерация по результатам запроса
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func GetBannerStorage(db *sql.DB, ids []int, filteredBanner *[]models.FilteredBanner) error {
	query := `
		SELECT
			id_banner, title_banner, text_banner, url_banner, banner_state, created_at, updated_at, tag_list, features_id
		FROM banners_storage
		WHERE id_banner = ANY($1);
	`

	rows, err := db.Query(query, pq.Array(ids))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var temp models.FilteredBanner
		err := rows.Scan(&temp.Id, &temp.Banner.Title, &temp.Banner.Text, &temp.Banner.Url, &temp.Flag, &temp.CreatedAt, &temp.UpdatedAt, &temp.TagIds, &temp.FeatureIds)
		if err != nil {
			return err
		}
		*filteredBanner = append(*filteredBanner, temp)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	fmt.Println(filteredBanner)
	return nil
}

func GetBannerHistoryStorage(db *sql.DB, id string, filteredBanner *[]models.HistoryBanner) error {
	query := `
		SELECT 
			id, features_id, tag_list, title_banner, text_banner, url_banner, banner_state, created_at, updated_at
		FROM history_banenrs 
		WHERE id_banner = $1;
	`

	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var temp models.HistoryBanner
		err := rows.Scan(&temp.Id, &temp.FeatureIds, &temp.TagIds, &temp.Banner.Title, &temp.Banner.Text, &temp.Banner.Url, &temp.Flag, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return err
		}
		*filteredBanner = append(*filteredBanner, temp)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
