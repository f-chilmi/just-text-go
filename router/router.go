package router

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/f-chilmi/just-text-go/controllers"
)

func Router() *mux.Router {
	fmt.Println("router called")
	router := mux.NewRouter()

	router.HandleFunc("/", controllers.HomeController).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", controllers.FindAll).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/{id}", controllers.FindById).Methods("GET", "OPTIONS")
	router.HandleFunc("/user", controllers.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT", "OPTIONS")

	router.HandleFunc("/new-msg", controllers.NewMsg).Methods("POST", "OPTIONS")

	// find user by phone
	router.HandleFunc("/phone/{phone}", controllers.FindByPhone).Methods("GET", "OPTIONS")

	// get rooms
	// by token
	router.HandleFunc("/room", controllers.ListRoom).Methods("GET", "OPTIONS")
	// by room id
	router.HandleFunc("/room/{id}", controllers.OpenRoom).Methods("GET", "OPTIONS")

	return router
}
