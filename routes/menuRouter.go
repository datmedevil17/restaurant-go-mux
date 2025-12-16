package routes

import (
	controller "github.com/datmedevil17/restaurant-management/controllers"
	"github.com/gorilla/mux"
)

func MenuRoutes(r *mux.Router) {
	r.HandleFunc("/menus", controller.GetMenus).Methods("GET")
	r.HandleFunc("/menus/{menu_id}", controller.GetMenu).Methods("GET")
	r.HandleFunc("/menus", controller.CreateMenu).Methods("POST")
	r.HandleFunc("/menus/{menu_id}", controller.UpdateMenu).Methods("PUT")
	r.HandleFunc("/menus/{menu_id}", controller.DeleteMenu).Methods("DELETE")

}
