package main

import (
	"log"
	"net/http"
	"os"

	"github.com/datmedevil17/restaurant-management/middlewares"
	"github.com/datmedevil17/restaurant-management/routes"
	"github.com/gorilla/mux"
)


// helper for client.DatabaseName("").collection(collectionName)


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.Use(middlewares.Logger)

	routes.UserRoutes(r)

	api := r.PathPrefix("/").Subrouter()
	api.Use(middlewares.Authentication)

	routes.UserProtectedRoutes(api)
	routes.FoodRoutes(api)
	routes.MenuRoutes(api)
	routes.OrderRoutes(api)
	routes.OrderItemRoutes(api)
	routes.TableRoutes(api)
	routes.InvoiceRoutes(api)

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
