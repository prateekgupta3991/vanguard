package handler

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
)

// statusHandler reports process health and service readiness.
func statusHandler(c *gin.Context) {
	respond(c, stdhttp.StatusOK, gin.H{"status": "ok"}, "VANGUARD_HEALTHY", "Vanguard is healthy")
}
