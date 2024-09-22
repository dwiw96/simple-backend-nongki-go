package main

import (
	"log"
	config "simple-backend-nongki-go/config"
	postgres "simple-backend-nongki-go/utils/driver/postgresql"
)

func main() {
	log.Println("-- run market maven v2 --")
	env := config.GetEnvConfig()
	pgPool := postgres.ConnectToPg(env)

	defer pgPool.Close()
}
