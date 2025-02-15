package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthAndSendCoins(t *testing.T) {
	authRequest1 := map[string]string{
		"username": "user1",
		"password": "password123",
	}
	authBody1, _ := json.Marshal(authRequest1)
	resp1, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authBody1))
	assert.NoError(t, err)
	defer resp1.Body.Close()

	var response1 map[string]string
	err = json.NewDecoder(resp1.Body).Decode(&response1)
	assert.NoError(t, err)
	token1 := response1["token"]

	authRequest2 := map[string]string{
		"username": "user2",
		"password": "password123",
	}
	authBody2, _ := json.Marshal(authRequest2)
	resp2, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authBody2))
	assert.NoError(t, err)
	defer resp2.Body.Close()

	var response2 map[string]string
	err = json.NewDecoder(resp2.Body).Decode(&response2)
	assert.NoError(t, err)

	sendCoinsRequest := map[string]interface{}{
		"toUser": "user2",
		"amount": 100,
	}
	sendCoinsBody, _ := json.Marshal(sendCoinsRequest)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/sendCoin", bytes.NewBuffer(sendCoinsBody))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token1)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
