package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthAndBuy(t *testing.T) {
	authRequest := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	authBody, _ := json.Marshal(authRequest)
	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	token := response["token"]

	buyRequest := map[string]string{
		"itemName": "t-shirt",
	}
	buyBody, _ := json.Marshal(buyRequest)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/buy/t-shirt", bytes.NewBuffer(buyBody))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
