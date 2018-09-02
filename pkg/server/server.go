package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
)

type serv struct {
	db     *sql.DB
	router *http.ServeMux
}

type record struct {
	Pk    string `json:"pk"`
	Score string `json:"score"`
}

func saveRecord(db *sql.DB, r record) error {
	_, err := db.Exec(
		`INSERT INTO records (
			pk,
			score,
			created_at
		) VALUES ($1, $2, now())`,
		r.Pk,
		r.Score,
	)

	return err
}

func (s *serv) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok")
	}

}

func (s *serv) handleRecords() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			s.handlePostRecordsRequest(w, req)
			return
		}

		handleError(
			w,
			http.StatusMethodNotAllowed,
			fmt.Errorf("Method [%s] Not Allowed", req.Method),
		)

	}
}

func (s *serv) handlePostRecordsRequest(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	status := http.StatusCreated
	for {
		var r record

		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			status = http.StatusBadRequest
		}

		if err := saveRecord(s.db, r); err != nil {
			log.Println(err)
			status = http.StatusUnprocessableEntity
		}
	}

	w.WriteHeader(status)
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	var msg string
	if err != nil {
		msg = fmt.Sprintf(`{"error": %q}`, err)
	}
	fmt.Sprintf("Error: %s", msg)
	fmt.Println(msg)
	http.Error(w, msg, statusCode)
}

func New(db *sql.DB) *http.ServeMux {
	s := serv{
		db:     db,
		router: http.NewServeMux(),
	}

	s.routes()

	return s.router
}
