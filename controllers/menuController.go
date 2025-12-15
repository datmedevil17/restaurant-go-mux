package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	database "github.com/datmedevil17/restaurant-management/databases"
	model "github.com/datmedevil17/restaurant-management/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func getMenu(menuId string) (*model.Menu, error) {
	fId, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		log.Fatal(err)
	}
	var menu model.Menu
	filter := bson.M{"_id": fId}
	err = menuCollection.FindOne(context.TODO(), filter).Decode(&menu)
	if err != nil {
		return nil, err
	}
	return &menu, err
}

func getMenus() ([]model.Menu, error) {
	var menus []model.Menu
	cursor, err := menuCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var menu model.Menu
		err := cursor.Decode(&menu)
		if err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func createMenu(menu model.Menu) (*model.Menu, error) {
	menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	result, err := menuCollection.InsertOne(context.TODO(), menu)
	if err != nil {
		return nil, err
	}
	menu.ID = result.InsertedID.(primitive.ObjectID)
	return &menu, nil
}

func updateMenu(menuId string, menu model.Menu) (*model.Menu, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fId, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": fId}
	var updateObj bson.D

	if menu.Start_Date != nil && menu.End_Date != nil {
		if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
			return nil, errors.New("kindly retype the time")
		}
		updateObj = append(updateObj,
			bson.E{Key: "start_date", Value: menu.Start_Date},
			bson.E{Key: "end_date", Value: menu.End_Date},
		)
	}

	if menu.Name != "" {
		updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
	}

	if menu.Category != "" {
		updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
	}

	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

	if len(updateObj) == 0 {
		return nil, errors.New("no fields to update")
	}

	_, err = menuCollection.UpdateOne(
		ctx,
		filter,
		bson.M{"$set": updateObj},
	)
	if err != nil {
		return nil, err
	}

	menu.ID = fId
	return &menu, nil
}

func deleteMenu(menuId string) error {
	fId, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": fId}
	_, err = menuCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

//getMenus
//getMenu
//createMenu
//updateMenu
//deleteMenu

func GetMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	params := mux.Vars(r)
	menuId := params["menu_id"]
	menu, err := getMenu(menuId)
	if err != nil {
		if err.Error() == "food not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Food not found"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menu)
}

func GetMenus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	menus, err := getMenus()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menus)
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var menu model.Menu
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Bad request"})
		return
	}
	createdMenu, err := createMenu(menu)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdMenu)
}

func DeleteMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	menuId := params["menu_id"]
	err := deleteMenu(menuId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Menu deleted successfully"})
}
func UpdateMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	menuId := params["menu_id"]

	var menu model.Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Bad request"})
		return
	}

	updatedMenu, err := updateMenu(menuId, menu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMenu)
}
