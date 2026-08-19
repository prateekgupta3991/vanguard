package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prateekgupta3991/vanguard/internal/handler/dto"
)

// respond writes a successful Vanguard response envelope.
func respond(c *gin.Context, status int, data any, code, description string) {
	c.JSON(status, dto.VanguardResponse{Data: data, Code: code, Description: description})
}

// respondError writes an unsuccessful Vanguard response envelope.
func respondError(c *gin.Context, status int, code, description, errorCode, errorDescription string) {
	c.JSON(status, dto.VanguardResponse{
		Code:        code,
		Description: description,
		Error: &dto.VanguardError{
			ErrorCode: errorCode,
			ErrorDesc: errorDescription,
		},
	})
}
