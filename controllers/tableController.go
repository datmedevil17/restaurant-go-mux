package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	model "github.com/datmedevil17/restaurant-management/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tables []model.Table
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := tableCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while listing tables"})
		return
	}
	defer result.Close(ctx)
	for result.Next(ctx) {
		var table model.Table
		if err := result.Decode(&table); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, table)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tables)
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	tableId := params["table_id"]
	var table model.Table

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while fetching the table"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(table)
}

func CreateTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var table model.Table

	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.ID = primitive.NewObjectID()
	table.Table_id = table.ID.Hex()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, insertErr := tableCollection.InsertOne(ctx, table)
	if insertErr != nil {
		msg := "table item was not created"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": msg})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	tableId := params["table_id"]
	var table model.Table

	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	var updateObj bson.D

	if table.Number_of_guests != nil {
		updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
	}

	if table.Table_number != nil {
		updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Table_number})
	}

	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.Updated_at})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"table_id": tableId}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "table update failed"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	tableId := params["table_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := tableCollection.DeleteOne(ctx, bson.M{"table_id": tableId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while deleting the table"})
		return
	}

	if result.DeletedCount < 1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "table with this ID not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Table deleted successfully"})
}
