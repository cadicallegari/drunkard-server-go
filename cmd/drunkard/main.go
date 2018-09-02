/*
	trying some archteture like
	or
	http://idiomaticgo.com/post/best-practice/server-project-layout/
	or
	https://github.com/bxcodec/go-clean-arch
	or
	https://github.com/golang-standards/project-layout
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"cadicallegari/drunkard/pkg/server"
)

var (
	// version is set at build time
	Version = "No version provided at build time"
)

func newDB() *sql.DB {
	// use env var or other configurable thing
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		"postgres",
		"postgrespasswd",
		"db",
		"5432",
		"drunkard",
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
	// port := helper.GetEnvVarString("SERVER_PORT")
	server := &http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     handler,
		ReadTimeout: time.Minute,
	}

	fmt.Printf("Starting server listening port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
