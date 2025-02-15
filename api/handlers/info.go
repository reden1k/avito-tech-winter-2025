package handlers

import (
	"avito-tech-winter-2025/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InfoHandler(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Только GET-запросы разрешены"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не найден"})
		return
	}
	token := authHeader[7:]

	infoResponse, err := services.HandleInfoRequest(token)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, infoResponse)
}
