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
	"log"
	"os"
)

type Env struct {
	DbHost		string
	DbUser		string
	DbPassword	string
	DbName		string
	Listen		string
}

func getEnv() Env {
	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("DB_PASSWORD is not set")
	}

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		host = "localhost"
	}

	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		user = "postgres"
	}

	name, ok := os.LookupEnv("DB_NAME")
	if !ok {
		name = "postgres"
	}

	listen, ok := os.LookupEnv("LISTEN")
	if !ok {
		listen = ":4280"
	}

	return Env{
		DbPassword: password,
		DbHost: host,
		DbUser: user,
		DbName: name,
		Listen: listen,
	}
}
