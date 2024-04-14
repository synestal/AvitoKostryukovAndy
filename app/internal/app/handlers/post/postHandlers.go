package post

import (
	model "awesomeProject/internal/app/handlers"
	"awesomeProject/internal/app/handlers/contextLogger"
	"awesomeProject/internal/app/handlers/jsonProducer"
	postsql "awesomeProject/internal/app/sqlDAO/post"
	help "awesomeProject/pkg/func"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"strings"
)

type PostFuncs interface {
	SetAdmin(*sql.DB, string, string) (int, error)
	CreateNewBanner(*sql.DB, string, string, string, []string, []string) (int, *postsql.BannerId, error)
	ChangeBanner(*sql.DB, string, string, string, string, []string, []string) (int, error)
	DeleteBanner(*sql.DB, string, string) (int, error)
	DeleteBannerByFeatureOrTag(*sql.DB, string, string, string, string, string) (int, error)
	ChangeBannersHistory(*sql.DB, string, string, string) (int, error)
}

type PostF struct {
	pos PostFuncs
}

func New(pos PostFuncs) *PostF {
	return &PostF{
		pos: pos,
	}
}

func (p *PostF) PostAdminStateHandler(db *sql.DB) gin.HandlerFunc { // POST: set_admin
	return func(c *gin.Context) {
		id := c.Query("id")
		state := c.Query("state")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("Post admin state handler", id+", "+state))
		logger.InfoContext(ctx, "Starting handler")

		if !help.IsNumeric(id) || state != "true" && state != "false" {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, err := p.pos.SetAdmin(db, id, state)
		jsonProducer.ProduceRequest(c, request, err, logger, ctx)
	}
}

func (p *PostF) CreateNewBannerHandler(db *sql.DB) gin.HandlerFunc {
	// http://localhost:8080/banner?admin_token=10&feature_id=15&tag_ids=22,12&content=notebooklovers,simpledescr,http://aboba.com&is_active=true
	return func(c *gin.Context) {
		token := c.Query("admin_token")
		feature := c.Query("feature_id")
		tags := c.Query("tag_ids")
		content := c.Query("content")
		active := c.Query("is_active")
		parsedContent := strings.Split(content, ",")
		parsedTags := strings.Split(tags, ",")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("Create new banner handler", token+", "+feature+", "+tags+", "+content+", "+active))
		logger.InfoContext(ctx, "Starting handler")

		if len(parsedTags) < 2 || !help.IsNumeric(token) || !help.IsNumeric(feature) || !help.AllNumeric(parsedTags) || content == "" || active != "true" && active != "false" {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, BannerId, err := p.pos.CreateNewBanner(db, token, feature, active, parsedContent, parsedTags)
		jsonProducer.ProduceJSON(c, request, BannerId, err, logger, ctx)
	}
}

func (p *PostF) ChangeBannerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("admin_token")
		bannerid := c.Query("id")
		feature := c.Query("feature_id")
		tags := c.Query("tag_ids")
		content := c.Query("content")
		active := c.Query("is_active")
		parsedContent := strings.Split(content, ",")
		parsedTags := strings.Split(tags, ",")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("Change banner handler", token+", "+bannerid+", "+feature+", "+tags+", "+content+", "+active))
		logger.InfoContext(ctx, "Starting handler")

		if len(parsedTags) < 2 || !help.IsNumeric(bannerid) || !help.IsNumeric(token) || !help.IsNumeric(feature) || !help.AllNumeric(parsedTags) || content == "" || active != "true" && active != "false" {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, err := p.pos.ChangeBanner(db, token, bannerid, feature, active, parsedContent, parsedTags)
		jsonProducer.ProduceRequest(c, request, err, logger, ctx)
	}
}

func (p *PostF) DeleteBannerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("admin_token")
		bannerid := c.Query("id")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("Delete banner handler", token+", "+bannerid))
		logger.InfoContext(ctx, "Starting handler")

		if !help.IsNumeric(bannerid) || !help.IsNumeric(token) {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, err := p.pos.DeleteBanner(db, token, bannerid)
		jsonProducer.ProduceRequest(c, request, err, logger, ctx)
	}
}

func (p *PostF) DeleteFeatureTagHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("admin_token")
		feature := c.Query("feature_id")
		tag := c.Query("tag_id")
		limit := c.Query("content")
		offset := c.Query("offset")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("Delete by feature and tag handler", token+", "+feature+", "+tag+", "+limit+", "+offset))
		logger.InfoContext(ctx, "Starting handler")

		if !help.IsNumeric(token) || feature == "" && tag == "" {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, err := p.pos.DeleteBannerByFeatureOrTag(db, token, feature, limit, offset, tag)
		if err != nil {

		}
		jsonProducer.ProduceRequest(c, request, err, logger, ctx)
	}
}

func (p *PostF) ChangeBannersHistoryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("admin_token")
		id := c.Query("id")
		number := c.Query("number")

		h := contextLogger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
		logger := slog.New(h)
		ctx := contextLogger.AppendCtx(context.Background(), slog.String("Change banners history handler", token+", "+id+", "+number))
		logger.InfoContext(ctx, "Starting handler")

		if !help.IsNumeric(token) || !help.IsNumeric(number) || !help.IsNumeric(id) {
			logger.ErrorContext(ctx, "Error in handler, err 400")
			var errJSONResponse model.ErrJSONResponse
			errJSONResponse.ErrorJSON = "400, некорректные данные"
			c.JSON(400, errJSONResponse)
			return
		}
		request, err := p.pos.ChangeBannersHistory(db, token, number, id)
		if err != nil {

		}
		jsonProducer.ProduceRequest(c, request, err, logger, ctx)
	}
}
