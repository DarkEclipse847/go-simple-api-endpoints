package dbconn_test

import (
	"database/sql"
	"net/http"
	"strconv"
	dbconn "testTask/dbConn"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestInsertTestUser(t *testing.T) {
	var test bool
	db, err := sql.Open("postgres", "user=postgres port=1337 password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer db.Close()
	err = dbconn.InsertTestUser(db, "9ca97cac-51ef-4107-b5ec-89ff2b571712", 1000.00)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	err = db.QueryRow("SELECT exists(SELECT 1 FROM wallets WHERE uuid = $1 LIMIT 1)", "9ca97cac-51ef-4107-b5ec-89ff2b571712").Scan(&test)
	if err != nil {
		log.Error().Msgf("Error occured while check uuid existence %v", err)
		return
	}
	//_, err = db.Exec("DELETE * FROM wallets WHERE uuid = '9ca97cac-51ef-4107-b5ec-89ff2b571712'")
	//if err != nil {
	//	log.Fatal().Msg(err.Error())
	//}
	assert.Equal(t, true, test)
}

func TestUpdateWalletBalance(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres port=1337 password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer db.Close()
	err = dbconn.UpdateWalletBalance(db, "9ca97cac-51ef-4107-b5ec-89ff2b571712", "WITHDRAW", 500.00, nil, http.Request{})
	if err != nil {
		log.Error().Msgf("Error occured while updating balance: %v", err)
		return
	}
	amount, err := dbconn.GetWalletBalance(db, "9ca97cac-51ef-4107-b5ec-89ff2b571712", nil, http.Request{})
	if err != nil {
		log.Error().Msgf("Error occured while getting balance: %v", err)
		return
	}
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Error().Msgf("Error occured while parsing amount: %v", err)
		return
	}
	assert.Equal(t, 500.00, amountFloat)
}

func TestGetWalletBalance(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres port=1337 password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer db.Close()

	// Assuming the test user was inserted in the previous test
	amount, err := dbconn.GetWalletBalance(db, "9ca97cac-51ef-4107-b5ec-89ff2b571712", nil, http.Request{})
	if err != nil {
		log.Error().Msgf("Error occured while getting balance: %v", err)
		return
	}
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Error().Msgf("Error occured while parsing amount: %v", err)
		return
	}

	assert.Equal(t, 500.00, amountFloat)
}
