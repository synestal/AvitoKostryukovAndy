package get_funcs

import (
	getsql "awesomeProject/internal/app/sqlDAO/get"
	"awesomeProject/internal/app/sqlDAO/models"
	casher "awesomeProject/internal/redis-casher"
	help "awesomeProject/pkg/func"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
)

type Get struct {
}

func New() *Get {
	return &Get{}
}

func (g *Get) GetAdminState(db *sql.DB, token string) (bool, bool, error) {
	avaliable, adminState, err := getsql.GetAdminState(db, token)
	return avaliable, adminState, err
}

func (g *Get) GetBannerFromCache(client *redis.Client, db *sql.DB, tagID, featureID string) (string, *models.Banner, error) {
	state, banner, err := casher.GetBannerFromCache(client, db, tagID, featureID)
	return state, banner, err
}

func (g *Get) GetBannerFromDB(db *sql.DB, tagID, featureID string) (*models.Banner, string, error) {
	banner, state, err := getsql.GetBannerFromDB(db, tagID, featureID)
	return banner, state, err
}

func (g *Get) GetBannerByFilter(db *sql.DB, token, feature, limit, offset, tag string) (int, []models.FilteredBanner, error) {
	filteredBanner := make([]models.FilteredBanner, 0, 1)
	avaliable, adminState, err := g.GetAdminState(db, token)
	if err != nil {
		return 400, filteredBanner, err
	}
	if !avaliable {
		return 401, filteredBanner, errors.New("unauthorized")
	}
	if !adminState {

		return 403, filteredBanner, errors.New("unauthorized")
	}

	if tag == "" {
		ids, err := getsql.GetBannerIdByFeature(db, feature, limit, offset)
		if err != nil {
			return 500, filteredBanner, err
		}
		err = getsql.GetBannerStorage(db, ids, &filteredBanner)
		if err != nil {
			return 500, filteredBanner, err
		}
	} else if feature == "" {
		ids, err := getsql.GetBannerIdByTag(db, tag, limit, offset)
		if err != nil {
			return 500, filteredBanner, err
		}
		err = getsql.GetBannerStorage(db, ids, &filteredBanner)
		if err != nil {
			return 500, filteredBanner, err
		}
	} else {
		var (
			idsFea []int
			idsTag []int
			errFea error
			errTag error
		)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			idsFea, errFea = getsql.GetBannerIdByFeature(db, feature, limit, offset)
		}()
		go func() {
			defer wg.Done()
			idsTag, errTag = getsql.GetBannerIdByTag(db, tag, limit, offset)
		}()
		wg.Wait()

		if errFea != nil {
			err = fmt.Errorf("error fetching IDs by feature: %w", errFea)
			return 500, nil, err
		}
		if errTag != nil {
			err = fmt.Errorf("error fetching IDs by tag: %w", errTag)
			return 500, nil, err
		}

		ids := help.Intersection(idsFea, idsTag)
		err = getsql.GetBannerStorage(db, ids, &filteredBanner)
		if err != nil {
			return 500, filteredBanner, err
		}
	}
	return 200, filteredBanner, nil
}

func (g *Get) GetBannersHistory(db *sql.DB, token, id string) (int, []models.HistoryBanner, error) {
	filteredBanner := make([]models.HistoryBanner, 0, 1)
	avaliable, adminState, err := getsql.GetAdminState(db, token)
	if err != nil {
		return 500, filteredBanner, errors.New("unauthorized")
	}
	if !avaliable {
		return 401, filteredBanner, errors.New("unauthorized")
	}
	if !adminState {
		return 403, filteredBanner, errors.New("unauthorized")
	}
	err = getsql.GetBannerHistoryStorage(db, id, &filteredBanner)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 404, filteredBanner, errors.New("banner history not found")
		}
		return 500, filteredBanner, err
	}
	return 200, filteredBanner, nil
}
