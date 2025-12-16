package routes

import (
	controller "github.com/datmedevil17/restaurant-management/controllers"
	"github.com/gorilla/mux"
)

func OrderItemRoutes(r *mux.Router) {
	r.HandleFunc("/order-items", controller.GetOrderItems).Methods("GET")
	r.HandleFunc("/order-items/{order_item_id}", controller.GetOrderItem).Methods("GET")
	r.HandleFunc("/order-items-order/{order_item_id}", controller.GetOrderItemsByOrder).Methods("GET")
	r.HandleFunc("/order-items", controller.CreateOrderItem).Methods("POST")

	r.HandleFunc("/order-items/{order_item_id}", controller.UpdateOrderItem).Methods("PUT")
	r.HandleFunc("/order-items/{order_item_id}", controller.DeleteOrderItem).Methods("DELETE")
}
