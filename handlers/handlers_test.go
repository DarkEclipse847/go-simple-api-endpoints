package handlers_test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testTask/models"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Run this only when starting, this is a smoke-test for GetBalance handler(asserts fixed starting value)
func TestGetBalanceHandler(t *testing.T) {
	var test models.WalletResp
	resp, err := http.Get("http://localhost:8000/api/v1/wallets/c071658f-7c70-48af-95c8-2a7cf46536f6")
	if err != nil {
		log.Error().Msgf("failed to call: %v", err)
	}
	json.NewDecoder(resp.Body).Decode(&test)

	testFloat, err := strconv.ParseFloat(test.Amount, 64)
	if err != nil {
		log.Error().Msgf("failed to parse: %v", err)
	}
	assert.Equal(t, testFloat, 15000.00)
	defer resp.Body.Close()
}
