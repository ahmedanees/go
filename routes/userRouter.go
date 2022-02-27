package routes

import (
	controller "abwaab/controllers"

	"github.com/gorilla/mux"
)

//UserRoutes function
func UserRoutes(incomingRoutes *mux.Router) {
	//incomingRoutes.Use(middleware.Authentication)
	incomingRoutes.HandleFunc("/api/user/creatTweet",controller.CreateTweet).Methods("POST")
	incomingRoutes.HandleFunc("/api/user/searchTweet",controller.SearchAndSaveTweet).Methods("POST")
	incomingRoutes.HandleFunc("/api/user/listTweet",controller.ListTweet).Methods("GET")
}
