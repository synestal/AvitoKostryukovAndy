package main

import (
	"awesomeProject/internal/app/sqlDAO/models"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", "localhost", 5432, "postgres", "Synesta17", "AvitoDb", "disable")
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("Error in test")
	}
	defer db.Close()
	TestGetUserBanner(db)
}

func TestGetUserBanner(db *sql.DB) {
	//"http://localhost:8080/user_banner?tag_id=8&feature_id=15&use_last_revision=true&admin_token=25"
	TagIds := []string{"1", "14", "22"}
	FeatureIds := "2"
	Banner := []string{"some_title", "some_text", "some_url"}
	Flag := "true"
	err := CreateNemBannerStorage(db, TagIds, FeatureIds, Flag, Banner)
	if err != nil {
		fmt.Println("Error in test")
	}

	val, err := GetBannerFromDB(db, "1", "2")
	if err != nil {
		fmt.Println(err)
	}
	if val.Url != Banner[2] || val.Text != Banner[1] || val.Title != Banner[0] {
		fmt.Println("Error in test")
	}
	fmt.Println("Good test")

}

func CreateNemBannerStorage(db *sql.DB, tags []string, feature, active string, content []string) error {
	query := `
		INSERT INTO banners_storage(tag_list, features_id, title_banner, text_banner, url_banner, banner_state)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.Exec(query, pq.Array(tags), feature, content[0], content[1], content[2], active)
	if err != nil {
		return err
	}

	return nil
}

func GetBannerFromDB(db *sql.DB, tagID, featureID string) (*models.Banner, error) {
	query := `
		SELECT title_banner, text_banner, url_banner
		FROM banners_storage
		WHERE features_id = $2 AND $1 = ANY(tag_list);
	`

	rows, err := db.Query(query, tagID, featureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banner models.Banner
	if rows.Next() {
		err := rows.Scan(&banner.Title, &banner.Text, &banner.Url)
		if err != nil {
			return nil, err
		}
	}

	return &banner, nil
}
