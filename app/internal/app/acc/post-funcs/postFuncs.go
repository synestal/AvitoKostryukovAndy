package post_funcs

import (
	getsql "awesomeProject/internal/app/sqlDAO/get"
	postsql "awesomeProject/internal/app/sqlDAO/post"
	"database/sql"
)

type Post struct {
}

func New() *Post {
	return &Post{}
}
func (p *Post) CreateNewBanner(db *sql.DB, token, feature, active string, content, tags []string) (int, *postsql.BannerId, error) {
	avaliable, adminState, err := getsql.GetAdminState(db, token)
	if err != nil {
		return 500, nil, err
	}
	if !avaliable {
		return 401, nil, err
	}
	if !adminState {
		return 403, nil, err
	}
	bannerId, err := postsql.CreateNemBannerStorage(db, tags, feature, active, content)
	if err != nil {
		return 500, nil, err
	}

	return 201, bannerId, nil
}

func (p *Post) ChangeBanner(db *sql.DB, token, bannerid, feature, active string, content, tags []string) (int, error) {
	avaliable, adminState, err := getsql.GetAdminState(db, token)
	if err != nil {
		return 500, err
	}
	if !avaliable {
		return 401, err
	}
	if !adminState {
		return 403, err
	}

	ans, err := postsql.UpdateBannersStorage(db, tags, feature, bannerid, active, content)
	if err != nil {
		return 500, err
	}
	if ans == "no banner" {
		return 404, nil
	}
	return 200, nil
}

func (p *Post) DeleteBanner(db *sql.DB, token, bannerid string) (int, error) {
	avaliable, adminState, err := getsql.GetAdminState(db, token)
	if err != nil {
		return 500, err
	}
	if !avaliable {
		return 401, err
	}
	if !adminState {
		return 403, err
	}

	ans, err := postsql.DeleterBanners(db, bannerid)
	if err != nil {
		return 500, err
	}
	if ans == "not found" {
		return 404, nil
	}
	return 204, nil
}

func (p *Post) DeleteBannerByFeatureOrTag(db *sql.DB, token, feature, limit, offset, tag string) (int, error) {
	avaliable, adminState, err := getsql.GetAdminState(db, token)
	if err != nil {
		return 500, err
	}
	if !avaliable {
		return 401, err
	}
	if !adminState {
		return 403, err
	}

	ids := make([]int, 0, 1)
	if tag == "" {
		ids, err = getsql.GetBannerIdByFeature(db, feature, limit, offset)
		if err != nil {
			return 400, err
		}
	} else {
		ids, err = getsql.GetBannerIdByTag(db, tag, limit, offset)
		if err != nil {
			return 400, err
		}
	}

	err = postsql.DeleterBannersPostponed(db, ids)
	if err != nil {
		return 500, err
	}

	return 200, nil
}

func (p *Post) ChangeBannersHistory(db *sql.DB, token, number, id string) (int, error) {
	avaliable, adminState, err := getsql.GetAdminState(db, token)
	if err != nil {
		return 500, err
	}
	if !avaliable {
		return 401, err
	}
	if !adminState {
		return 403, err
	}

	ans, err := postsql.ChangeHistoryBannersStorage(db, number, id)
	if err != nil {
		return 500, err
	}
	if ans == "NULL" {
		return 404, nil
	}
	return 200, nil
}

func (p *Post) SetAdmin(db *sql.DB, id, state string) (int, error) {
	err := postsql.SetadminState(db, id, state)
	if err != nil {
		return 500, err
	}
	return 200, nil
}
