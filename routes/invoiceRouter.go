package routes

import (
	controller "github.com/datmedevil17/restaurant-management/controllers"
	"github.com/gorilla/mux"
)

func InvoiceRoutes(r *mux.Router){
	r.HandleFunc("/invoices", controller.GetInvoices()).Methods("GET")
	r.HandleFunc("/invoices/:invoice_id", controller.GetInvoice()).Methods("GET")
	r.HandleFunc("/invoices", controller.CreateInvoice()).Methods("POST")
	r.HandleFunc("/invoices/:invoice_id", controller.UpdateInvoice()).Methods("PUT")
	r.HandleFunc("/invoices/:invoice_id", controller.DeleteInvoice()).Methods("DELETE")

}