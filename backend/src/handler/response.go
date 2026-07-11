package handler

import (
	"go-perjalanan-dinas/dto"

	"github.com/gin-gonic/gin"
)

func writeSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, dto.WebResponse{Success: true, Message: message, Data: data})
}

func writeSuccessMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, dto.WebResponse{Success: true, Message: message})
}

func writeError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, dto.WebResponse{Success: false, Message: message})
}
