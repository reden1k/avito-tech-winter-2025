package handlers

import (
	"avito-tech-winter-2025/dto"
	"avito-tech-winter-2025/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SendCoinsHandler(c *gin.Context) {
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

	var sendCoinsRequest dto.SendCoinsRequest
	if err := c.ShouldBindJSON(&sendCoinsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования запроса"})
		return
	}

	sendCoinsResponse, err := services.HandleSendCoinsRequest(token, sendCoinsRequest)
	if err != nil {
		log.Printf("Ошибка при переводе монет: %s", err.Message)
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	log.Printf("Успешный перевод %d монет пользователю %s", sendCoinsRequest.Amount, sendCoinsRequest.ReceiverUsername)

	c.JSON(http.StatusOK, sendCoinsResponse)
}
