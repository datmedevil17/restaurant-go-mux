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

type InvoiceViewFormat struct {
	Invoice_id       string      `json:"invoice_id"`
	Order_id         string      `json:"order_id"`
	Payment_method   *string     `json:"payment_method" validate:"required,eq=CARD|eq=CASH|eq="`
	Payment_status   *string     `json:"payment_status" validate:"required,eq=PENDING|eq=PAID|eq="`
	Payment_due_date time.Time   `json:"payment_due_date"`
	Payment_due      interface{} `json:"payment_due"`
	Table_number     interface{} `json:"table_number"`
	Order_details    interface{} `json:"order_details"`
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var invoices []model.Invoice
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := invoiceCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while listing invoice items"})
		return
	}
	defer result.Close(ctx)
	for result.Next(ctx) {
		var invoice model.Invoice
		if err := result.Decode(&invoice); err != nil {
			log.Fatal(err)
		}
		invoices = append(invoices, invoice)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoices)
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	invoiceId := params["invoice_id"]
	var invoice model.Invoice

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while listing invoice item"})
		return
	}

	var invoiceView InvoiceViewFormat
	invoiceView.Invoice_id = invoice.Invoice_id
	invoiceView.Order_id = invoice.Order_id
	invoiceView.Payment_method = invoice.Payment_method
	invoiceView.Payment_status = invoice.Payment_status
	invoiceView.Payment_due_date = invoice.Payment_due_date

	// Get Order Details
	var order model.Order
	err = orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while finding order"})
		return
	}

	// Get Table Details
	var table model.Table
	if order.Table_id != nil {
		err = tableCollection.FindOne(ctx, bson.M{"table_id": *order.Table_id}).Decode(&table)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "error occured while finding table"})
			return
		}
		invoiceView.Table_number = table.Table_number
	}

	// Get Order Items and Calculate Payment Due
	var orderItems []model.OrderItem
	cursor, err := orderItemCollection.Find(ctx, bson.M{"order_id": invoice.Order_id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while finding order items"})
		return
	}
	defer cursor.Close(ctx)

	var paymentDue float64 = 0
	for cursor.Next(ctx) {
		var orderItem model.OrderItem
		if err := cursor.Decode(&orderItem); err != nil {
			log.Fatal(err)
		}
		if orderItem.Unit_price != nil {
			paymentDue += *orderItem.Unit_price
		}
		orderItems = append(orderItems, orderItem)
	}
	invoiceView.Payment_due = paymentDue
	invoiceView.Order_details = orderItems

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceView)
}

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var invoice model.Invoice
	var order model.Order

	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
	if err != nil {
		msg := "message: Order was not found"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.Invoice_id = invoice.ID.Hex()

	result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
	if insertErr != nil {
		msg := "invoice item was not created"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": msg})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	invoiceId := params["invoice_id"]
	var invoice model.Invoice

	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while decoding the request body"})
		return
	}

	var updateObj bson.D

	if invoice.Payment_method != nil {
		updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.Payment_method})
	}
	if invoice.Payment_status != nil {
		updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.Payment_status})
	}

	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.Updated_at})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"invoice_id": invoiceId}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := invoiceCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "invoice item update failed"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	invoiceId := params["invoice_id"]

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := invoiceCollection.DeleteOne(ctx, bson.M{"invoice_id": invoiceId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "error occured while deleting the invoice item"})
		return
	}

	if result.DeletedCount < 1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "invoice with this ID not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invoice deleted successfully"})
}
