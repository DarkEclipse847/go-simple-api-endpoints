package dbconn

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func CreateTable(db *sql.DB) error {

	// [Why numeric type instead of money?](https://wiki.postgresql.org/wiki/Don%27t_Do_This#Don.27t_use_money)
	// Why numeric instead of float? Because float has some bugs with rounding and mathematical operations, cf.(https://www.reddit.com/r/learnprogramming/comments/1gv27ya/floating_point_errors/)

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS wallets (id SERIAL PRIMARY KEY, uuid VARCHAR(36) UNIQUE, amount NUMERIC(10, 2) DEFAULT 0.00)")
	if err != nil {
		log.Error().Msgf("Error occured while creating table %v", err)
		return err
	}
	return nil
}

func InsertTestUser(db *sql.DB, uuid string, amount float32) error {
	var exists bool
	// Ограничиваем запрос (LIMIT 1) чтобы поиск выполнялся быстрее
	err := db.QueryRow("SELECT exists(SELECT 1 FROM wallets WHERE uuid = $1 LIMIT 1)", uuid).Scan(&exists)
	if err != nil {
		log.Error().Msgf("Error occured while checking uuid existence %v", err)
		return err
	}
	// UUIDs generated randomly using [this](https://www.uuidgenerator.net/) cause there is no api endpoint to generate it(login/auth etc.)
	if !exists {
		_, err = db.Exec("INSERT INTO wallets (uuid, amount) VALUES($1, $2)", uuid, amount)
		if err != nil {
			log.Info().Msgf("Error occured while inserting test value %v", err)
			return err
		}
	}

	return nil
}

func UpdateWalletBalance(db *sql.DB, uuid string, operation string, amount float64, w http.ResponseWriter, r http.Request) error {
	// Using trasactionsan isolation level to fullfil ACID principles
	tx, err := db.BeginTx(r.Context(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		log.Error().Msgf("Transaction crashed at start: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	switch operation {
	case "DEPOSIT":
		_, err := tx.Exec("UPDATE wallets SET amount = amount + $1 WHERE uuid = $2", amount, uuid)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error occured while updating value %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	case "WITHDRAW":
		_, err := tx.Exec("UPDATE wallets SET amount = amount - $1 WHERE uuid = $2", amount, uuid)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error occured while updating test value %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	default:
		log.Error().Msgf("Unsufficient 'OPERATION_TYPE' value: %v", operation)
	}
	err = tx.Commit()
	if err != nil {
		log.Error().Msgf("Error occured while commiting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func GetWalletBalance(db *sql.DB, uuid string, w http.ResponseWriter, r http.Request) (string, error) {
	tx, err := db.BeginTx(r.Context(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		log.Error().Msgf("Transaction crashed at start: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	var balance string
	err = tx.QueryRow("SELECT amount FROM wallets WHERE uuid = $1 LIMIT 1", uuid).Scan(&balance)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error occured while updating test value %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		log.Error().Msgf("Error occured while commiting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	return balance, nil
}
