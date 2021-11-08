package api

import (
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) Server {
	return Server{db}
}

func (server Server) index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello\n")
}

func (server Server) Serve(listen string) error {
	http.HandleFunc("/", server.index)
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listen, nil)
}
