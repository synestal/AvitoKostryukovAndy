package models

type Banner struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type FilteredBanner struct {
	Id         string `json:"banner_id"`
	TagIds     string `json:"tag_ids"`
	FeatureIds string `json:"feature_id"`
	Banner     Banner `json:"content"`
	Flag       string `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type HistoryBanner struct {
	Id         string `json:"change_id"`
	TagIds     string `json:"tag_ids"`
	FeatureIds string `json:"feature_id"`
	Banner     Banner `json:"content"`
	Flag       string `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
