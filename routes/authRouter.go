package routes

import (
	controller "abwaab/controllers"

	"github.com/gorilla/mux"
)

//UserRoutes function
func AuthRoutes(incomingRoutes *mux.Router) {
	incomingRoutes.HandleFunc("/api/user/login",controller.UserLogin).Methods("POST")
	incomingRoutes.HandleFunc("/api/user/signup",controller.UserSignup).Methods("POST")
}
