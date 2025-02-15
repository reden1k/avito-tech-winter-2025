package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	authRequest := map[string]string{
		"username": "testuser323232fs",
		"password": "password123",
	}

	authBody, _ := json.Marshal(authRequest)

	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(authBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])
}
