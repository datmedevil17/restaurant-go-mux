package routes

import (
	controller "github.com/datmedevil17/restaurant-management/controllers"
	"github.com/gorilla/mux"
)

func OrderRoutes(r *mux.Router){
	r.HandleFunc("/orders", controller.GetOrders).Methods("GET")
	r.HandleFunc("/orders/:order_id", controller.GetOrder).Methods("GET")
	r.HandleFunc("/orders", controller.CreateOrder).Methods("POST")
	r.HandleFunc("/orders/:order_id", controller.UpdateOrder).Methods("PUT")
	r.HandleFunc("/orders/:order_id", controller.DeleteOrder).Methods("DELETE")
}