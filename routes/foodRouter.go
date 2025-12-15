package routes

import (
	controller "github.com/datmedevil17/restaurant-management/controllers"
	"github.com/gorilla/mux"
)

func FoodRoutes(r *mux.Router) {
	r.HandleFunc("/foods", controller.GetFoods).Methods("GET")
	r.HandleFunc("/foods/:food_id", controller.GetFood).Methods("GET")
	r.HandleFunc("/foods", controller.CreateFood).Methods("POST")
	r.HandleFunc("/foods/:food_id", controller.UpdateFood).Methods("PATCH")
	r.HandleFunc("/foods/:food_id", controller.DeleteFood).Methods("DELETE")
}
