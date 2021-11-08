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
