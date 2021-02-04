package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"ad-service/internal/ad"
	"ad-service/internal/config"
)

var flagConfig = flag.String("config", "./config/config.yml", "path to the config file")

func main() {
	flag.Parse()

	// load config
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatalf("failed connect to postgres: %s", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(db),
	}

	log.Printf("server is running at %v", address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}


func buildHandler(db *sql.DB) http.Handler {
	router := mux.NewRouter()

	a := ad.New(db)
	a.RegisterHandlers(router)

	return router
}
