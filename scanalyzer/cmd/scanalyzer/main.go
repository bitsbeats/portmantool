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

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bitsbeats/portmantool/scanalyzer/internal/api"
	"github.com/bitsbeats/portmantool/scanalyzer/internal/database"
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
