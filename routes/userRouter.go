package routes

import (
	controller "github.com/datmedevil17/restaurant-management/controllers"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/users", controller.GetUsers).Methods("GET")
	r.HandleFunc("/users/:user_id", controller.GetUser).Methods("GET")
	r.HandleFunc("/users/signup", controller.SignUp).Methods("POST")
	r.HandleFunc("/users/login", controller.Login).Methods("POST")
	// r.HandleFunc("/users/:user_id", controller.UpdateUser).Methods("PUT")
	// r.HandleFunc("/users/:user_id", controller.DeleteUser).Methods("DELETE")

}
