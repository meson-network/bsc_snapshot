package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type MSG_RESP_HEALTH struct {
	Unixtime int64 `json:"unixtime"`
}

// @Summary      /api/health
// @Description  health check
// @Tags         health
// @Produce      json
// @Success      200 {object} MSG_RESP_HEALTH "server unix time"
// @Router       /api/health [get]
func healthCheck(ctx echo.Context) error {
	// o := 0
	// v := 5 / o
	// return errors.New(strconv.Itoa(v))
	return ctx.JSON(http.StatusOK, &MSG_RESP_HEALTH{Unixtime: time.Now().Unix()})
}

func configHealth(httpServer *echo.Echo) {
	//health
	httpServer.GET("/api/health", healthCheck)
}
