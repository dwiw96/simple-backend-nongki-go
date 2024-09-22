package main

import (
	"context"
	"log"
	server "simple-backend-nongki-go/api"
	config "simple-backend-nongki-go/config"
	factory "simple-backend-nongki-go/factory"
	postgres "simple-backend-nongki-go/utils/driver/postgresql"
	password "simple-backend-nongki-go/utils/password"
	"time"
)

func main() {
	log.Println("-- run market maven v2 --")
	env := config.GetEnvConfig()
	pgPool := postgres.ConnectToPg(env)
	defer pgPool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	password.JwtInit(pgPool, ctx)

	router := server.SetupRouter()

	factory.InitFactory(router, pgPool, ctx)

	server.StartServer(env.SERVER_PORT, router)
}
