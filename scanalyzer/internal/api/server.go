package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitsbeats/portmantool/scanalyzer/internal/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) Server {
	return Server{db}
}

func (server Server) Serve(listen string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/diff", server.diffExpected)
	mux.HandleFunc("/diff/", server.diffTwo)
	mux.HandleFunc("/expected", server.expected)
	mux.HandleFunc("/hello", server.hello)
	mux.HandleFunc("/run", server.run)
	mux.HandleFunc("/run/", server.run)
	mux.HandleFunc("/scans", server.getScans)
	mux.HandleFunc("/scans/", server.deleteScans)
	mux.HandleFunc("/scan", badRequest("id required"))
	mux.HandleFunc("/scan/", server.scan)

	http.Handle("/v1", http.NotFoundHandler())
	http.Handle("/v1/", http.StripPrefix("/v1", mux))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listen, nil)
}

func (server Server) diffExpected(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		diff, err := database.DiffExpected(server.db)
		if err != nil {
			serverError(w, r, err)
			return
		}

		toJSON(w, r, diff)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) diffTwo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(501)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) expected(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		expectedState, err := database.Expected(server.db)
		if err != nil {
			serverError(w, r, err)
			return
		}

		toJSON(w, r, expectedState)
	case "PATCH":
		if r.Header.Get("Content-Type") != "application/json" {
			clientError(w, r, "data must be sent as json")
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			serverError(w, r, err)
			return
		}

		var expectedState []database.ExpectedState
		err = json.Unmarshal(data, &expectedState)
		if err != nil {
			clientError(w, r, err)
			return
		}

		err = server.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&expectedState).Error
		if err != nil {
			serverError(w, r, err)
			return
		}

		w.WriteHeader(204)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello\n")
}

func (server Server) run(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.WriteHeader(501)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) getScans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		scans, err := database.Scans(server.db)
		if err != nil {
			serverError(w, r, err)
			return
		}

		toJSON(w, r, scans)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) deleteScans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		w.WriteHeader(501)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) scan(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := stripPrefix(r, "/scan/")
		if id == "" {
			clientError(w, r, "id must not be empty")
			return
		}

		timestamp, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			clientError(w, r, err)
			return
		}

		scan, err := database.StateAt(server.db, time.Unix(timestamp, 0))
		if err != nil {
			serverError(w, r, err)
			return
		}

		toJSON(w, r, scan)
	default:
		w.WriteHeader(405)
	}
}

func badRequest(message string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientError(w, r, message)
	}
}

func clientError(w http.ResponseWriter, r *http.Request, v interface{}) {
	httpError(w, r, 400, v)
}

func serverError(w http.ResponseWriter, r *http.Request, err error) {
	httpError(w, r, 500, err)
}

func httpError(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	log.Printf("%s %s - %v", r.Method, r.URL.RequestURI(), v)
	w.WriteHeader(status)
	io.WriteString(w, fmt.Sprint(v))
}

func toJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		serverError(w, r, err)
		return
	}

	w.Write(data)
}

func stripPrefix(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, prefix)
}
