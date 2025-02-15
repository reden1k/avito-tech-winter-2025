package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthBuyAndSendCoinsWithInfo(t *testing.T) {
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
	token2 := response2["token"]

	buyRequest1 := map[string]string{
		"itemName": "t-shirt",
	}
	buyBody1, _ := json.Marshal(buyRequest1)
	req1, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/buy/t-shirt", bytes.NewBuffer(buyBody1))
	assert.NoError(t, err)
	req1.Header.Set("Authorization", "Bearer "+token1)

	client := &http.Client{}
	resp1, err = client.Do(req1)
	assert.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	buyRequest2 := map[string]string{
		"itemName": "cup",
	}
	buyBody2, _ := json.Marshal(buyRequest2)
	req2, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/buy/cup", bytes.NewBuffer(buyBody2))
	assert.NoError(t, err)
	req2.Header.Set("Authorization", "Bearer "+token2)

	resp2, err = client.Do(req2)
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	sendCoinsRequest := map[string]interface{}{
		"toUser": "user2",
		"amount": 100,
	}
	sendCoinsBody, _ := json.Marshal(sendCoinsRequest)
	req3, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/sendCoin", bytes.NewBuffer(sendCoinsBody))
	assert.NoError(t, err)
	req3.Header.Set("Authorization", "Bearer "+token1)

	resp3, err := client.Do(req3)
	assert.NoError(t, err)
	defer resp3.Body.Close()
	assert.Equal(t, http.StatusOK, resp3.StatusCode)

	req4, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/info", nil)
	assert.NoError(t, err)
	req4.Header.Set("Authorization", "Bearer "+token1)

	resp4, err := client.Do(req4)
	assert.NoError(t, err)
	defer resp4.Body.Close()
	assert.Equal(t, http.StatusOK, resp4.StatusCode)

	req5, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/info", nil)
	assert.NoError(t, err)
	req5.Header.Set("Authorization", "Bearer "+token2)

	resp5, err := client.Do(req5)
	assert.NoError(t, err)
	defer resp5.Body.Close()
	assert.Equal(t, http.StatusOK, resp5.StatusCode)
}
