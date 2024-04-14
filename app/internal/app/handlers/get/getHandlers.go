package get

import (
	model "awesomeProject/internal/app/handlers"
	"awesomeProject/internal/app/handlers/contextLogger"
	"awesomeProject/internal/app/handlers/jsonProducer"
	"awesomeProject/internal/app/sqlDAO/models"
	help "awesomeProject/pkg/func"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log/slog"
	"os"
)

type GetFuncs interface {
	GetAdminState(*sql.DB, string) (bool, bool, error)                                                       // Done
	GetBannerFromCache(*redis.Client, *sql.DB, string, string) (string, *models.Banner, error)               //Done
	GetBannerFromDB(*sql.DB, string, string) (*models.Banner, string, error)                                 // Done
	GetBannerByFilter(*sql.DB, string, string, string, string, string) (int, []models.FilteredBanner, error) // In progress
	GetBannersHistory(*sql.DB, string, string) (int, []models.HistoryBanner, error)
}

type GetF struct {
	gt GetFuncs
}

func New(gt GetFuncs) *GetF {
	return &GetF{
		gt: gt,
	}
}

func (p *GetF) GetBannerHandler(db *sql.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID := c.Query("tag_id")
		featureID := c.Query("feature_id")
		useLast := c.Query("use_last_revision")
		token := c.Query("admin_token")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("get banner handler", tagID+", "+featureID+", "+useLast+", "+token))
		logger.InfoContext(ctx, "Starting handler")

		if !help.IsNumeric(tagID) || !help.IsNumeric(featureID) || useLast != "false" && useLast != "true" || !help.IsNumeric(token) {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}

		useLastRevision := useLast == "true"
		var adminState bool
		var err error
		avaliable, adminState, err := p.gt.GetAdminState(db, token)
		if err != nil {
			logger.ErrorContext(ctx, "Error in handler, err 500", slog.String("error", err.Error()))
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "500, внутренняя ошибка сервера"
			c.JSON(500, errJSONResponse)
			return
		}

		// Проверка кэша
		if !useLastRevision {
			state, banner, err := p.gt.GetBannerFromCache(redisClient, db, tagID, featureID)
			if err != nil {
				logger.ErrorContext(ctx, "Error in handler, err 500", slog.String("error", err.Error()))
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "500, внутренняя ошибка сервера"
				c.JSON(500, errJSONResponse)
				return
			}
			if !avaliable {
				logger.ErrorContext(ctx, "Error in handler, err 401")
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "401, пользователь не авторизован"
				c.JSON(401, errJSONResponse)
				return
			}
			if state == "false" && !adminState {
				logger.ErrorContext(ctx, "Error in handler, err 403")
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "403, пользователь не имеет доступа"
				c.JSON(403, errJSONResponse)
				return
			}
			if banner.Title == "" && banner.Text == "" && banner.Url == "" {
				logger.ErrorContext(ctx, "Error in handler, err 404")
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "404, баннер не найден"
				c.JSON(404, errJSONResponse)
				return
			}
			logger.InfoContext(ctx, "Success handler")
			c.JSON(200, banner)
		}

		if useLastRevision {
			var state string
			banner, state, err := p.gt.GetBannerFromDB(db, tagID, featureID)
			if err != nil {
				logger.ErrorContext(ctx, "Error in handler, err 500", slog.String("error", err.Error()))
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "500, внутренняя ошибка сервера"
				c.JSON(500, errJSONResponse)
				return
			}
			if !avaliable {
				logger.ErrorContext(ctx, "Error in handler, err 401")
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "401, пользователь не авторизован"
				c.JSON(401, errJSONResponse)
				return
			}
			if state == "false" && !adminState {
				logger.ErrorContext(ctx, "Error in handler, err 403")
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "403, пользователь не имеет доступа"
				c.JSON(403, errJSONResponse)
				return
			}
			if banner.Title == "" && banner.Text == "" && banner.Url == "" {
				logger.ErrorContext(ctx, "Error in handler, err 404")
				var errJSONResponse model.ErrJSONResponse
				errJSONResponse.ErrorJSON = "404, баннер не найден"
				c.JSON(404, errJSONResponse)
				return
			}
			logger.InfoContext(ctx, "Success handler")
			c.JSON(200, banner)
		}
	}
}

func (p *GetF) GetBannerByFilterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("admin_token")
		feature := c.Query("feature_id")
		tag := c.Query("tag_id")
		limit := c.Query("content")
		offset := c.Query("offset")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("get banner by filter handler", token+", "+feature+", "+tag+", "+limit+", "+offset))
		logger.InfoContext(ctx, "Starting handler")

		if !help.IsNumeric(token) || !help.IsNumeric(feature) && help.IsNumeric(tag) {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, bannerList, err := p.gt.GetBannerByFilter(db, token, feature, limit, offset, tag)
		if len(bannerList) == 0 {
			request = 404
		}
		jsonProducer.ProduceJSON(c, request, bannerList, err, logger, ctx)
	}
}

func (p *GetF) GetBannersHistoryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("admin_token")
		id := c.Query("id")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("get banner history handler", token+", "+id))
		logger.InfoContext(ctx, "Starting handler")

		if !help.IsNumeric(token) || !help.IsNumeric(id) {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, bannerList, err := p.gt.GetBannersHistory(db, token, id)
		if len(bannerList) == 0 {
			request = 404
		}
		jsonProducer.ProduceJSON(c, request, bannerList, err, logger, ctx)
	}
}
