package routes

import (
	controller "github.com/datmedevil17/restaurant-management/controllers"
	"github.com/gorilla/mux"
)

func TableRoutes(r *mux.Router) {
	r.HandleFunc("/tables", controller.GetTables).Methods("GET")
	r.HandleFunc("/tables/{table_id}", controller.GetTable).Methods("GET")
	r.HandleFunc("/tables", controller.CreateTable).Methods("POST")
	r.HandleFunc("/tables/{table_id}", controller.UpdateTable).Methods("PUT")
	r.HandleFunc("/tables/{table_id}", controller.DeleteTable).Methods("DELETE")

}
