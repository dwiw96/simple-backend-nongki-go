package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter() *httprouter.Router {
	log.Println("<- SetupRouter()")

	router := httprouter.New()

	log.Println("-> SetupRouter()")
	return router
}

func StartServer(port string, router *httprouter.Router) {
	log.Println("start server at localhost", port)
	log.Fatal(http.ListenAndServe(port, router))
}
