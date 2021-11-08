package main

import (
	"context"
	"log"

	"github.com/bitsbeats/portmantool/scanalyzer/internal/api"
	"github.com/bitsbeats/portmantool/scanalyzer/internal/database"
	"github.com/bitsbeats/portmantool/scanalyzer/internal/importer"
	"github.com/bitsbeats/portmantool/scanalyzer/internal/metrics"
)

func main() {
	env := getEnv()

	db, err := database.Connect(env.DbHost, env.DbUser, env.DbPassword, env.DbName)
	if err != nil {
		log.Fatal(err)
	}

	err = metrics.RegisterMetrics()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := importer.NewImporter(db)
	err = i.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(db)
	log.Fatal(server.Serve(env.Listen))
}
