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
	log.Fatal(http.ListenAndServe(port, router))
}
