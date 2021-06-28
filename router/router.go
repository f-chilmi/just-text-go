package router

import (
	"github.com/gorilla/mux"

	"github.com/f-chilmi/just-text-go/controllers"
	"github.com/f-chilmi/just-text-go/middlewares"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// authentications
	router.HandleFunc("/register", controllers.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", controllers.Login).Methods("POST", "OPTIONS")

	// users
	router.HandleFunc("/", middlewares.SetMiddlewareAuth(controllers.HomeController)).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", middlewares.SetMiddlewareAuth(controllers.FindAll)).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/{id}", middlewares.SetMiddlewareAuth(controllers.FindById)).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/{id}", middlewares.SetMiddlewareAuth(controllers.UpdateUser)).Methods("PUT", "OPTIONS")

	// find user by phone
	router.HandleFunc("/phone/{phone}", middlewares.SetMiddlewareAuth(controllers.FindRoomByPhone)).Methods("GET", "OPTIONS")

	// get rooms
	// by token
	router.HandleFunc("/room", middlewares.SetMiddlewareAuth(controllers.ListRoom)).Methods("GET", "OPTIONS")
	// by room id
	router.HandleFunc("/room/{id}", middlewares.SetMiddlewareAuth(controllers.OpenRoom)).Methods("GET", "OPTIONS")

	// send message
	router.HandleFunc("/msg/{id}", middlewares.SetMiddlewareAuth(controllers.SendMsg)).Methods("POST", "OPTIONS")

	return router
}
