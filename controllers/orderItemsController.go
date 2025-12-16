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

type OrderItemPack struct {
	Table_id    *string           `json:"table_id"`
	Order_items []model.OrderItem `json:"order_items"`
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "order_item")

func GetOrderItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orderItems []model.OrderItem
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderItemCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while listing order items"})
		return
	}
	defer result.Close(ctx)
	for result.Next(ctx) {
		var orderItem model.OrderItem
		if err := result.Decode(&orderItem); err != nil {
			log.Fatal(err)
		}
		orderItems = append(orderItems, orderItem)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderItems)
}

func GetOrderItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderItemId := params["order_item_id"]
	var orderItem model.OrderItem

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while listing order item"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderItem)
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localField", Value: "food_id"}, {Key: "foreignField", Value: "food_id"}, {Key: "as", Value: "food"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "order"}, {Key: "localField", Value: "order_id"}, {Key: "foreignField", Value: "order_id"}, {Key: "as", Value: "order"}}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$order"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "table"}, {Key: "localField", Value: "order.table_id"}, {Key: "foreignField", Value: "table_id"}, {Key: "as", Value: "table"}}}}
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$table"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "amount", Value: "$food.price"},
			{Key: "total_count", Value: 1},
			{Key: "food_name", Value: "$food.name"},
			{Key: "food_image", Value: "$food.food_image"},
			{Key: "table_number", Value: "$table.table_number"},
			{Key: "table_id", Value: "$table.table_id"},
			{Key: "order_id", Value: "$order.order_id"},
			{Key: "price", Value: "$food.price"},
			{Key: "quantity", Value: 1},
		}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "order_id", Value: "$order_id"}, {Key: "table_id", Value: "$table_id"}, {Key: "table_number", Value: "$table_number"}}}, {Key: "order_items", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}

	projectStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: 1},
		}}}

	var allStages mongo.Pipeline

	allStages = append(allStages, matchStage)
	allStages = append(allStages, lookupStage)
	allStages = append(allStages, unwindStage)
	allStages = append(allStages, lookupOrderStage)
	allStages = append(allStages, unwindOrderStage)
	allStages = append(allStages, lookupTableStage)
	allStages = append(allStages, unwindTableStage)
	allStages = append(allStages, projectStage)
	allStages = append(allStages, groupStage)
	allStages = append(allStages, projectStage2)

	result, err := orderItemCollection.Aggregate(ctx, allStages)

	if err != nil {
		return nil, err
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		return nil, err
	}

	return OrderItems, err
}

func GetOrderItemsByOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderId := params["order_id"]

	allOrderItems, err := ItemsByOrder(orderId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while listing order items by order ID"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allOrderItems)
}

func CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orderItemPack OrderItemPack
	var order model.Order

	if err := json.NewDecoder(r.Body).Decode(&orderItemPack); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItemsToBeInserted := []interface{}{}
	order.Table_id = orderItemPack.Table_id
	order_id := OrderItemOrderCreator(order)

	for _, orderItem := range orderItemPack.Order_items {
		orderItem.Order_id = order_id

		orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.ID = primitive.NewObjectID()
		orderItem.Order_item_id = orderItem.ID.Hex()

		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}

	insertedOrderItems, err := orderItemCollection.InsertMany(context.TODO(), orderItemsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insertedOrderItems)
}

func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orderItem model.OrderItem
	params := mux.Vars(r)
	orderItemId := params["order_item_id"]

	if err := json.NewDecoder(r.Body).Decode(&orderItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	var updateObj bson.D

	if orderItem.Unit_price != nil {
		updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.Unit_price})
	}

	if orderItem.Quantity != nil {
		updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
	}

	if orderItem.Food_id != nil {
		updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.Food_id})
	}

	orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

	filter := bson.M{"order_item_id": orderItemId}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderItemCollection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "order item update failed"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	orderItemId := params["order_item_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderItemCollection.DeleteOne(ctx, bson.M{"order_item_id": orderItemId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while deleting the order item"})
		return
	}

	if result.DeletedCount < 1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "order item with this ID not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order item deleted successfully"})
}
