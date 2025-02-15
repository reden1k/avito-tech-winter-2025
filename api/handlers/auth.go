package handlers

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthHandler(c *gin.Context) {
	var req dto.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга JSON"})
		return
	}

	authResponse, authErr := services.HandleAuthRequest(req)
	if authErr != nil {
		c.JSON(authErr.StatusCode, gin.H{"error": authErr.Message, "code": authErr.Code})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}
