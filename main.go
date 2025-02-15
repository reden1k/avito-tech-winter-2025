package main

import (
	"avito-tech-winter-2025/api/handlers"
	"avito-tech-winter-2025/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	r.POST("/api/auth", handlers.AuthHandler)
	r.GET("/api/info", handlers.InfoHandler)
	r.POST("/api/buy/:itemName", handlers.BuyHandler)
	r.POST("/api/sendCoin", handlers.SendCoinsHandler)

	log.Println("Сервер запущен на порту 8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}

	authReq := Request{
		Method: "POST",
		URL:    "http://localhost:8080/api/auth",
		Body: map[string]string{
			"username": "testuser",
			"password": "pass",
		},
	}

	resp, err := SendRequest(authReq)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Статус ответа:", resp.Status)
}

type Request struct {
	Method  string
	URL     string
	Body    interface{}
	Headers map[string]string
	Token   string
}

func SendRequest(reqData Request) (*http.Response, error) {
	var req *http.Request
	var err error

	var bodyReader *strings.Reader
	if reqData.Body != nil {
		var jsonData []byte
		jsonData, err = json.Marshal(reqData.Body)
		if err != nil {
			return nil, fmt.Errorf("ошибка сериализации тела запроса: %v", err)
		}
		bodyReader = strings.NewReader(string(jsonData))
	} else {
		bodyReader = strings.NewReader("")
	}

	req, err = http.NewRequest(reqData.Method, reqData.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	for key, value := range reqData.Headers {
		req.Header.Set(key, value)
	}

	if reqData.Token != "" {
		req.Header.Set("Authorization", "Bearer "+reqData.Token)
	}

	if reqData.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %v", err)
	}

	return resp, nil
}
