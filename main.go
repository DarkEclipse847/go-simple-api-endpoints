package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	dbconn "testTask/dbConn"
	"testTask/handlers"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	//init logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	//db init
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer db.Close()

	err = dbconn.CreateTable(db)
	if err != nil {
		log.Fatal().Msg(err.Error())

	}
	err = dbconn.InsertTestUser(db, "c071658f-7c70-48af-95c8-2a7cf46536f6", 15000.0)
	if err != nil {
		log.Fatal().Msg(err.Error())

	}
	err = dbconn.InsertTestUser(db, "768aa6f7-f304-4ed3-be2d-aba7df964534", 10000.0)
	if err != nil {
		log.Fatal().Msg(err.Error())

	}

	router := mux.NewRouter()

	mainRouter := router.Methods(http.MethodPost, http.MethodGet).Subrouter()

	mainRouter.HandleFunc("/api/v1/wallet", handlers.WalletOperationHandler(db)).Methods("POST")
	mainRouter.HandleFunc("/api/v1/wallets/{uuid}", handlers.GetBalanceHandler(db)).Methods("GET")

	probesRouter := router.Methods(http.MethodGet).Subrouter()
	probesRouter.HandleFunc("/probes/readiness",
		func(rw http.ResponseWriter, r *http.Request) {
			_, err := rw.Write([]byte("OK"))
			if err != nil {
				log.Error().Msgf("Error while writing the data to an HTTP reply with err=%s", err)
				return
			}
		})

	probesRouter.HandleFunc("/probes/liveness", func(rw http.ResponseWriter, r *http.Request) {

		//check if we can access DB
		connStr := os.Getenv("DATABASE_URL")

		db, err := sql.Open("postgres", connStr)
		log.Info().Msgf("Successful db connect: %v", db)

		if err != nil {
			log.Error().Msgf("Error while connection to DB with err=%s", err)
			return
		}
	})

	port := os.Getenv("APPLICATION_PORT")
	if len(port) == 0 {
		log.Fatal().Msgf("APPLICATION_PORT env doesnot not set")
	}
	srvPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal().Msgf("cannot cast APPLICATION_PORT env to integer")
	}

	log.Info().Msgf("starting the server on port :%s", port)

	//http.Server instance
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", srvPort),
		Handler:      router,
		TLSConfig:    nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Msgf("Starting server on port %d", srvPort)

		err := s.ListenAndServe()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	//trap os.Signal and gracefully shutdown the server
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)
	sig := <-sigCh
	log.Info().Msgf("Graceful shutdown with signal %s \n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
