// +build integration

// https://medium.com/@benbjohnson/structuring-tests-in-go-46ddee7a25c
// https://segment.com/blog/5-advanced-testing-techniques-in-go/
// https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831

package server_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"cadicallegari/drunkard/pkg/server"
)

func assert(tb testing.TB, condition bool, msg string) {
	if !condition {
		tb.Error(msg)
	}
}

func ok(tb testing.TB, err error, msg string) {
	if err != nil {
		tb.Errorf("Error not expected: %s\n", err)
	}
}

func equals(tb testing.TB, exp, act interface{}, msg string) {
	if exp != act {
		tb.Errorf("%s, not equals: %q, but got: %q", msg, exp, act)
	}
}

func newDB(t *testing.T) *sql.DB {
	// use env var or other configurable thing
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		"postgres",
		"postgrespasswd",
		"db",
		"5432",
		"drunkard-dev",
	)

	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestHandleHello(t *testing.T) {
	srv := server.New(newDB(t))

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hello", nil)
	ok(t, err, "creating about request")

	srv.ServeHTTP(res, req)
	equals(t, res.Code, http.StatusOK, "response code")
}

func TestHandleListRestaurants(t *testing.T) {
	srv := server.New(newDB(t))

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/restaurants", nil)
	ok(t, err, "creating about request")

	srv.ServeHTTP(res, req)
	equals(t, res.Code, http.StatusOK, "response code")
}
