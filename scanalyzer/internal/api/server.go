// Copyright 2020-2022 Thomann Bits & Beats GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bitsbeats/portmantool/scanalyzer/internal/database"
	"github.com/bitsbeats/portmantool/scanalyzer/internal/importer"
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

func (server Server) ListenAndServe(listen string, ctx context.Context, wg *sync.WaitGroup) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/diff", server.diffExpected)
	mux.HandleFunc("/diff/", server.diffTwo)
	mux.HandleFunc("/expected", server.expected)
	mux.HandleFunc("/hello", server.hello)
	mux.HandleFunc("/run", server.run)
	mux.HandleFunc("/run/", server.run)
	mux.HandleFunc("/scans", server.scans)
	mux.HandleFunc("/scans/", server.pruneScans)
	mux.HandleFunc("/scan", server.scan)
	mux.HandleFunc("/scan/", server.scan)

	http.Handle("/v1", http.NotFoundHandler())
	http.Handle("/v1/", http.StripPrefix("/v1", mux))
	http.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr: listen,
	}

	go func() {
		wg.Add(1)
		defer wg.Done()

		<-ctx.Done()

		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Print(err)
		}
	}()

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
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
		ids := strings.Split(stripPrefix(r, "/diff/"), "/")
		if ids[0] == "" {
			clientError(w, r, "id1 must not be empty")
			return
		}
		if len(ids) == 2 && ids[1] == "" {
			clientError(w, r, "id2 must not be empty when given")
			return
		}
		if len(ids) > 2 {
			clientError(w, r, "too many path components")
			return
		}

		id1, err := strconv.ParseInt(ids[0], 10, 64)
		if err != nil {
			clientError(w, r, err)
			return
		}

		switch len(ids) {
		case 1:
			diff, err := database.DiffOne(server.db, time.Unix(id1, 0))
			if err != nil {
				serverError(w, r, err)
				return
			}

			toJSON(w, r, diff)
		case 2:
			id2, err := strconv.ParseInt(ids[1], 10, 64)
			if err != nil {
				clientError(w, r, err)
				return
			}

			diff, err := database.DiffTwo(server.db, time.Unix(id1, 0), time.Unix(id2, 0))
			if err != nil {
				serverError(w, r, err)
				return
			}

			toJSON(w, r, diff)
		}
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

func (server Server) scans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		scans, err := database.Scans(server.db)
		if err != nil {
			serverError(w, r, err)
			return
		}

		toJSON(w, r, scans)
	case "DELETE":
		err := database.Prune(server.db, time.Now())
		if err != nil {
			serverError(w, r, err)
			return
		}

		w.WriteHeader(204)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) pruneScans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		keep := stripPrefix(r, "/scans/")
		if keep == "" {
			clientError(w, r, "keep must not be empty when given")
			return
		}

		timestamp, err := strconv.ParseInt(keep, 10, 64)
		if err != nil {
			clientError(w, r, err)
			return
		}

		err = database.Prune(server.db, time.Unix(timestamp, 0))
		if err != nil {
			serverError(w, r, err)
			return
		}

		w.WriteHeader(204)
	default:
		w.WriteHeader(405)
	}
}

func (server Server) scan(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := stripPrefix(r, "/scan/")
		if !strings.HasPrefix(r.URL.Path, "/scan/") || id == "" {
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
	case "POST":
		if r.Header.Get("Content-Type") != "application/xml" && r.Header.Get("Content-Type") != "text/xml" {
			clientError(w, r, "scan report must be sent as xml")
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			serverError(w, r, err)
			return
		}

		err = importer.Import(server.db, data)
		if err != nil {
			serverError(w, r, err)
			return
		}

		w.WriteHeader(204)
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
