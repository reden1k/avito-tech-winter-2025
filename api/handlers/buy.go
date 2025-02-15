package handlers

import (
	"avito-tech-winter-2025/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BuyHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Только POST-запросы разрешены"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не найден"})
		return
	}
	token := authHeader[7:]

	itemName := c.Param("itemName")

	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя товара не указано"})
		return
	}

	buyResponse, err := services.HandleBuyRequest(token, itemName)
	if err != nil {
		log.Printf("Ошибка при обработке запроса: %s", err.Message)
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	log.Printf("Успешная покупка товара: %s", itemName)

	c.JSON(http.StatusOK, buyResponse)
}
