package main

import (
	"log"
	"net/http"
	"os"

	database "github.com/datmedevil17/restaurant-management/databases"
	"github.com/datmedevil17/restaurant-management/middlewares"
	"github.com/datmedevil17/restaurant-management/routes"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// helper for client.DatabaseName("").collection(collectionName)
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.Use(middlewares.Logger)
	r.Use(middlewares.Authentication)

	routes.UserRoutes(r)
	routes.FoodRoutes(r)
	routes.MenuRoutes(r)
	routes.OrderRoutes(r)
	routes.OrderItemRoutes(r)
	routes.TableRoutes(r)
	routes.InvoiceRoutes(r)

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
