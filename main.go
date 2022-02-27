package main

import (
	routes "abwaab/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)



func main(){

	log.Println("Starting the application")

	// Setup Routes for application
	router:= mux.NewRouter()
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	

	log.Fatal(http.ListenAndServe(":8080", router))
}
