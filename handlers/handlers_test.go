package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testTask/handlers"
	"testTask/models"
	"testing"

	_ "github.com/lib/pq"
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
func TestWalletOperationHandler(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres port=1337 password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer db.Close()
	reqBody := models.WalletReq{
		UUID:      "c071658f-7c70-48af-95c8-2a7cf46536f6",
		Operation: "DEPOSIT",
		Amount:    "1000",
	}
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/api/v1/wallet/", bytes.NewBuffer(body))
	if err != nil {
		log.Error().Msgf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.WalletOperationHandler(db)) // Assuming db is initialized
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

}
