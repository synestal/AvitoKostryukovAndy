package jsonProducer

import (
	model "awesomeProject/internal/app/handlers"
	"awesomeProject/internal/app/sqlDAO/models"
	"awesomeProject/internal/app/sqlDAO/post"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type MyConstraint interface {
	[]models.FilteredBanner |
		[]models.HistoryBanner |
		*post.BannerId
}

func ProduceJSON[T MyConstraint](c *gin.Context, request int, bannerList T, err error, logger *slog.Logger, ctx context.Context) {
	var errJSONResponse model.ErrJSONResponse
	errP := "normal"
	if err != nil {
		errP = err.Error()
	}
	if request == 500 {
		logger.ErrorContext(ctx, "Error in get func, err 500", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "500, внутренняя ошибка сервера"
		c.JSON(request, errJSONResponse)
		return
	} else if request == 403 {
		logger.ErrorContext(ctx, "Error in get func, err 403", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "403, пользователь не имеет доступа"
		c.JSON(request, errJSONResponse)
		return
	} else if request == 401 {
		logger.ErrorContext(ctx, "Error in get func, err 401", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "401, пользователь не авторизован"
		c.JSON(request, errJSONResponse)
		return
	} else if request == 404 {
		logger.ErrorContext(ctx, "Error in get func, err 404", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "404, баннер для тега не найден"
		c.JSON(request, errJSONResponse)
		return
	}
	logger.InfoContext(ctx, "Success handler")
	c.JSON(request, bannerList)
}

func ProduceRequest(c *gin.Context, request int, err error, logger *slog.Logger, ctx context.Context) {
	var errJSONResponse model.ErrJSONResponse
	errP := "normal"
	if err != nil {
		errP = err.Error()
	}
	if request == 500 {
		logger.ErrorContext(ctx, "Error in get func, err 500", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "500, внутренняя ошибка сервера"
	} else if request == 403 {
		logger.ErrorContext(ctx, "Error in get func, err 403", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "403, пользователь не имеет доступа"
	} else if request == 401 {
		logger.ErrorContext(ctx, "Error in get func, err 401", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "401, пользователь не авторизован"
	} else if request == 404 {
		logger.ErrorContext(ctx, "Error in get func, err 404", slog.String("error", errP))
		errJSONResponse.ErrorJSON = "404, баннер для тега не найден"
	} else if request == 204 {
		logger.InfoContext(ctx, "Success handler")
		errJSONResponse.ErrorJSON = "204, Баннер успешно удален"
	} else {
		logger.InfoContext(ctx, "Success handler")
		errJSONResponse.ErrorJSON = "200, ОК"
	}

	c.JSON(request, errJSONResponse)
}
