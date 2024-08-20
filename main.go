package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Meal struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type Order struct {
	ID     string `json:"id"`
	MealID string `json:"meal_id"`
	Amount int    `json:"amount"`
}

func main() {
	http.HandleFunc("/available-meals", availableMealsHandler)
	http.HandleFunc("/orders", ordersHandler)

	// Serve the static images
	fs := http.FileServer(http.Dir("./public/images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func availableMealsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonFile, err := os.Open("available-meals.json")
	if err != nil {
		http.Error(w, "Could not open available-meals.json", http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var meals []Meal
	json.Unmarshal(byteValue, &meals)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		jsonFile, err := os.Open("orders.json")
		if err != nil {
			http.Error(w, "Could not open orders.json", http.StatusInternalServerError)
			return
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var orders []Order
		json.Unmarshal(byteValue, &orders)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)

	case http.MethodPost:
		var newOrder Order
		err := json.NewDecoder(r.Body).Decode(&newOrder)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		jsonFile, err := os.Open("orders.json")
		if err != nil {
			http.Error(w, "Could not open orders.json", http.StatusInternalServerError)
			return
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var orders []Order
		json.Unmarshal(byteValue, &orders)

		orders = append(orders, newOrder)

		ordersJson, _ := json.Marshal(orders)
		err = ioutil.WriteFile("orders.json", ordersJson, 0644)
		if err != nil {
			http.Error(w, "Could not write to orders.json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newOrder)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
