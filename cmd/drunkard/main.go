package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cadicallegari/drunkard/pkg/server"
)

var (
	// version is set at build time
	Version = "No version provided at build time"
)

func newDB() *sql.DB {
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_SERVICE_HOST"),
		os.Getenv("POSTGRES_SERVICE_PORT"),
		os.Getenv("POSTGRES_DB_NAME"),
	)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {

	handler := server.New(newDB())

	port := "8080"
	server := &http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     handler,
		ReadTimeout: time.Minute,
	}

	fmt.Printf("Starting server listening port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
