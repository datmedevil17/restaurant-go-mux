package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	database "github.com/datmedevil17/restaurant-management/databases"
	model "github.com/datmedevil17/restaurant-management/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func getFood(foodId string) (*model.Food, error) {
	fId, err := primitive.ObjectIDFromHex(foodId)
	if err != nil {
		return nil, err
	}

	var food model.Food
	filter := bson.M{"_id": fId}
	err = foodCollection.FindOne(context.TODO(), filter).Decode(&food)
	if err != nil {
		return nil, err
	}
	return &food, err
}

func getFoods() ([]model.Food, error) {
	var foods []model.Food
	cursor, err := foodCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var food model.Food
		err := cursor.Decode(&food)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}

func createFood(food model.Food) (*model.Food, error) {
	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	result, err := foodCollection.InsertOne(context.TODO(), food)
	if err != nil {
		return nil, err
	}
	food.ID = result.InsertedID.(primitive.ObjectID)
	return &food, nil
}

func updateFood(foodId string, food model.Food) (*model.Food, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	fId, err := primitive.ObjectIDFromHex(foodId)
	if err != nil {
		return nil, err
	}

	var updateObj bson.D

	if food.Name != "" {
		updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
	}

	if food.Price != 0 {
		updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
	}

	if food.Food_image != "" {
		updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.Food_image})
	}

	if food.Menu_id != nil {
		updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.Menu_id})
	}

	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.Updated_at})

	filter := bson.M{"_id": fId}

	_, err = foodCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}})
	if err != nil {
		return nil, err
	}

	food.ID = fId
	return &food, nil
}

func deleteFood(foodId string) error {
	fId, err := primitive.ObjectIDFromHex(foodId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": fId}
	_, err = foodCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

//getFoods
//getFood
//createFood
//updateFood
//deleteFood

func GetFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	foodId := params["food_id"]
	food, err := getFood(foodId)
	if err != nil {
		if err.Error() == "food not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Food not found"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(food)
}

func GetFoods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	foods, err := getFoods()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foods)
}

func CreateFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var food model.Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Bad request"})
		return
	}
	createdFood, err := createFood(food)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdFood)
}

func DeleteFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	foodId := params["food_id"]
	err := deleteFood(foodId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Food deleted successfully"})
}

func UpdateFood(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	foodId := params["food_id"]
	var food model.Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Bad request"})
		return
	}
	updatedFood, err := updateFood(foodId, food)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedFood)
}
