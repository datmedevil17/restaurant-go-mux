package controller

import (
	"context"
	"encoding/json"
	"log"

	"net/http"
	"time"

	database "github.com/datmedevil17/restaurant-management/databases"
	model "github.com/datmedevil17/restaurant-management/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orders []model.Order
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while listing order items"})
		return
	}
	defer result.Close(ctx)
	for result.Next(ctx) {
		var order model.Order
		if err := result.Decode(&order); err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderId := params["order_id"]
	var order model.Order

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while fetching the order item"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order model.Order
	var table model.Table

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	if order.Table_id != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		tableID := *order.Table_id
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableID}).Decode(&table)
		if err != nil {
			msg := "message: Table was not found"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(msg)
			return
		}
	}

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, insertErr := orderCollection.InsertOne(ctx, order)
	if insertErr != nil {
		msg := "order item was not created"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": msg})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderId := params["order_id"]
	var order model.Order
	var table model.Table

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	var updateObj bson.D

	if order.Table_id != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		tableID := *order.Table_id
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableID}).Decode(&table)
		if err != nil {
			msg := "message: Table was not found"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(msg)
			return
		}
		updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.Table_id})
	}

	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.Updated_at})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"order_id": orderId}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "order item update failed"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderId := params["order_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderCollection.DeleteOne(ctx, bson.M{"order_id": orderId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while deleting the order item"})
		return
	}

	if result.DeletedCount < 1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "order with this ID not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order deleted successfully"})
}

func OrderItemOrderCreator(order model.Order) string {
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderCollection.InsertOne(ctx, order)
	return order.Order_id
}
