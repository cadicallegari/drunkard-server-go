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
	"os"
	"strings"
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
		tb.Errorf("%s, not equals: expecting %v, but got: %v", msg, exp, act)
	}
}

func newDB(t *testing.T) *sql.DB {
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
		t.Fatal(err)
	}

	return db
}

func setup(t *testing.T) (*http.ServeMux, func()) {
	db := newDB(t)
	return server.New(db), func() {
		db.Exec("TRUNCATE TABLE records CASCADE")
	}
}

func TestShouldBeHealth(t *testing.T) {
	srv := server.New(newDB(t))

	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/healthz", nil)
	ok(t, err, "creating healthz request")

	srv.ServeHTTP(res, req)
	equals(t, res.Code, http.StatusOK, "response code")
}

func TestHandleNewRecordProperly(t *testing.T) {
	srv, teardown := setup(t)
	defer teardown()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/records",
		strings.NewReader(`{"pk": "50", "score": "98"}`),
	)

	srv.ServeHTTP(res, req)
	equals(t, http.StatusCreated, res.Code, "status code")
}

func TestHandleBunchOfNewRecordsProperly(t *testing.T) {
	srv, teardown := setup(t)
	defer teardown()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/records",
		strings.NewReader(`
			{"pk": "51", "score": "98"}
			{"pk": "52", "score": "98"}
			{"pk": "53", "score": "98"}
			{"pk": "54", "score": "98"}
		`),
	)

	srv.ServeHTTP(res, req)
	equals(t, http.StatusCreated, res.Code, "status code")
}

func TestShouldReturnNotProcessedWhenPKAlreadExists(t *testing.T) {
	srv, teardown := setup(t)
	defer teardown()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/records",
		strings.NewReader(`
			{"pk": "same", "score": "98"}
			{"pk": "same", "score": "98"}
		`),
	)

	srv.ServeHTTP(res, req)
	equals(t, http.StatusUnprocessableEntity, res.Code, "status code")
}
