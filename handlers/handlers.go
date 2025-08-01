package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	dbconn "testTask/dbConn"
	"testTask/models"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func WalletOperationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Debug().Msgf("Request body: %v\n", r.Body)

		var wallet models.WalletReq
		json.NewDecoder(r.Body).Decode(&wallet)
		amountFloat, err := strconv.ParseFloat(wallet.Amount, 32)
		if err != nil {
			amountInt, err := strconv.ParseInt(wallet.Amount, 10, 32)
			//maybe there is no need to do this (was error in input while testing)
			if err != nil {
				log.Error().Msgf("Error occured while parsing amount: %v", err)
				http.Error(w, "Invalid amount", http.StatusInternalServerError)
				return
			}
			amountFloat = float64(amountInt)
		}

		err = dbconn.UpdateWalletBalance(db, wallet.UUID, wallet.Operation, float64(amountFloat), w, *r)
		if err != nil {
			log.Error().Msgf("Error occured while parsing amount: %v", err)
			http.Error(w, "Invalid amount", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Succsessfully updated balance"))
	}
}

func GetBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp models.WalletResp

		vars := mux.Vars(r)
		walletId := vars["uuid"]
		amount, err := dbconn.GetWalletBalance(db, walletId, w, *r)
		if err != nil {
			log.Error().Msgf("Error occured while getting balance(db connection issue?): %v", err)
			return
		}
		resp.Amount = amount
		res, err := json.Marshal(resp)
		if err != nil {
			log.Error().Msgf("client: Could not read json body: %v", err)
			http.Error(w, "Unable to read response", http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}
}
