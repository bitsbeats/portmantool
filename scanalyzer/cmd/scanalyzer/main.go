package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

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

	err = metrics.UpdateFromDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		wg.Add(1)
		defer wg.Done()

		<-ctx.Done()
		stop()
	}()
	defer stop() // technically unnecessary

	i := importer.NewImporter(db)
	err = i.Run(ctx, &wg)
	if err != nil {
		stop()
		wg.Wait()
		log.Fatal(err)
	}

	server := api.NewServer(db)
	err = server.ListenAndServe(env.Listen, ctx, &wg)
	if err != nil {
		stop()
		wg.Wait()
		log.Fatal(err)
	}

	wg.Wait()
	log.Print("bye")
}
