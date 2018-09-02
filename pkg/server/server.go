/*
trying to follow some advices from
https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831
*/

package server

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"net/http"
)

type serv struct {
	db     *sql.DB
	router *http.ServeMux
	// email  EmailSender
}

func (s *serv) handleHello() http.HandlerFunc {
	msg := `
    <html><body>
        <h1>Hello!</h1>
        <h1>%s</h1>
        <form method='POST' action='/hello'>
            <h2>What would you like me to say?</h2>
            <input name="message" type="text" >
            <input type="submit" value="Submit">
        </form>
    </body></html>`

	return func(w http.ResponseWriter, r *http.Request) {
		content := ""

		switch r.Method {
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			content = r.FormValue("message")
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, fmt.Sprintf(msg, content))

	}

}

func (s *serv) handleRestaurants() http.HandlerFunc {
	// fp := path.Join("templates", "index.html")

	// tmpl, err := template.ParseFiles(fp)
	// if err != nil {
	//     http.Error(w, err.Error(), http.StatusInternalServerError)
	//     return
	// }

	// if err := tmpl.Execute(w, book); err != nil {
	//  http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	tmpl := template.Must(template.ParseFiles("/templates/index.tmpl.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		// content := ""

		switch r.Method {
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// content = r.FormValue("message")
		}

		// w.WriteHeader(http.StatusOK)
		// w.Header().Set("Content-Type", "text/html")
		// io.WriteString(w, fmt.Sprintf(msg, content))

		bla := struct {
			Msg string
		}{
			Msg: "ma oeee",
		}

		if err := tmpl.Execute(w, bla); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

}

func New(db *sql.DB) *http.ServeMux {
	s := serv{
		db:     db,
		router: http.NewServeMux(),
	}

	s.routes()

	return s.router
}
