package http

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/vanguard/internal/http/handler/dto"
)

func NewHandler() stdhttp.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.HandleMethodNotAllowed = true
	router.GET("/healthz", statusHandler)
	router.GET("/readyz", statusHandler)
	return router
}

func statusHandler(c *gin.Context) {
	c.JSON(stdhttp.StatusOK, dto.VanguardResponse{
		Data:        gin.H{"status": "ok"},
		Code:        "VANGUARD_HEALTHY",
		Description: "Vanguard is healthy",
		Error:       nil,
	})
}
