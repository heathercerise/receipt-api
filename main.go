package main

// @title Receipt API
// @description A simple receipt processor
// @version 1.0.0

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Receipt structure
type Receipt struct {
	ID           string
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Item structure to be contained in receipts
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Response when request for points
type PointsResponse struct {
	Points int64 `json:"points"`
}

// Response when creating a new receipt
type IDResponse struct {
	ID string `json:"id"`
}

// Holds all receipts in program, normally would be a database
var receipts []Receipt

// Method to find a receipt given an ID in request
func GetReceiptByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		fmt.Println("ID isn't in the params")
	}

	for _, receipt := range receipts {
		if receipt.ID == id {
			// If found, calculate points and return JSON points object
			points := GetReceiptPoints(receipt)
			pointsStruct := PointsResponse{Points: points}
			json.NewEncoder(w).Encode(pointsStruct)
			return
		}
	}

	// If receipt not found, return 404 error
	http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
}

// Calculates receipts points with given instructions
func GetReceiptPoints(receipt Receipt) int64 {
	// One point for every alphanumeric character in retailer name
	retailer := receipt.Retailer
	points := GetAlphanumeric(retailer)

	// Points for total cost
	costStr := receipt.Total
	points += GetTotalCostPoints(costStr)

	// 5 points for every two items
	points += GetItemPoints(receipt)

	//iff generated using a large language model, 5 points if total is greater than 10.0
	// I assume this is a safeguard against using AI so skipping this?

	// 6 points if day in purchase date is odd
	dateString := receipt.PurchaseDate
	points += GetDatePoints(dateString)

	// 10 points if purchase between 2-4pm
	timeString := receipt.PurchaseTime
	points += GetTimePoints(timeString)

	return points
}

/*
	Below are various helper functions to help calculate receipt points
*/

// Returns number of alphanumeric characters
func GetAlphanumeric(str string) int64 {
	var total int64
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			total += 1
		}
	}
	return total
}

// Returns points given for total cost
func GetTotalCostPoints(costStr string) int64 {
	var points int64
	costFloat, err := strconv.ParseFloat(costStr, 64)
	if err == nil {
		costInt := int(math.Round(costFloat * 100))
		// 50 points if total is round dollar amount
		if costInt%100 == 0 {
			points += 50
		}
		// 25 points if total is multiple of .25
		if costInt%25 == 0 {
			points += 25
		}
	}

	return points
}

// Returns points for the items
func GetItemPoints(receipt Receipt) int64 {
	var numItems int
	var points int64

	// Five points for every two items
	for _, item := range receipt.Items {
		numItems += 1
		if numItems%2 == 0 {
			points += 5
		}

		// Trim item description and add points if multiple of 3
		desc := item.ShortDescription
		trimmed := strings.TrimSpace(desc)
		length := len(trimmed)
		if length%3 == 0 {
			itemPriceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				itemPriceFloat *= .2
				itemPrice := int(math.Ceil(itemPriceFloat))
				points += int64(itemPrice)
			}

		}
	}

	return points
}

// 6 points if bought on an odd day
func GetDatePoints(dateString string) int64 {
	var points int64
	day, err := strconv.Atoi(dateString[len(dateString)-2:])
	if err == nil {
		if day%2 == 1 {
			points = 6
		}
	}
	return points
}

// 10 points if after 2 and before 4 (14:00:00 to 15:59:59 is my assumption here)
func GetTimePoints(timeString string) int64 {
	var points int64
	time, err := strconv.Atoi(timeString[:2])
	if err == nil {
		if time >= 14 && time < 16 {
			points = 10
		}
	}
	return points
}

// Method to create a receipt with receipt json in the request; ensures valid receipt
func CreateReceipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var receipt Receipt
	_ = json.NewDecoder(r.Body).Decode(&receipt)

	// Validate fields
	// Description
	validReceipt := CheckValidDescription(receipt.Retailer)
	if !validReceipt {
		// Invalid receipt, set 400 error
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	// PurchaseDate and PurchaseTime
	validReceipt = CheckValidTime(receipt.PurchaseDate, receipt.PurchaseTime)
	if !validReceipt {
		// Invalid receipt, set 400 error
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	// Checks valid regex for both price and description
	validReceipt = CheckItemsValidity(receipt)
	if !validReceipt {
		// Invalid receipt, set 400 error
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	// Total cost
	validReceipt = CheckPriceValidity(receipt.Total)
	if !validReceipt {
		// Invalid receipt, set 400 error
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	} else {
		// Generate a unique ID for each receipt
		receipt.ID = GenerateID()
		receipts = append(receipts, receipt)

		// Return the ID JSON object of the created Receipt
		idStruct := IDResponse{ID: receipt.ID}
		json.NewEncoder(w).Encode(idStruct)
	}

}

/*
	Below are helper functions for creating and validating a receipt
*/

// Checks validity of description
func CheckValidDescription(str string) bool {
	valid, err := regexp.MatchString("^[\\w\\s\\-&]+$", str)
	if !valid || err != nil {
		fmt.Println("Retailer wrong format")
		return false
	}
	return true
}

// Checks validity of price
func CheckPriceValidity(str string) bool {
	valid, err := regexp.MatchString("\\d+\\.\\d{2}$", str)
	if !valid || err != nil {
		fmt.Println("Issue with total cost format")
		return false
	}

	return true
}

// Checks validity of date and time formatting
func CheckValidTime(dateString string, timeString string) bool {
	// PurchaseDate
	_, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Println("Invalid date format")
		return false
	}

	// PurchaseTime
	_, err = time.Parse("15:04", timeString)
	if err != nil {
		fmt.Println("Invalid time format")
		return false
	}

	// Valid time and date
	return true
}

// Checks validity of items
func CheckItemsValidity(receipt Receipt) bool {
	// Must be at least one item
	if len(receipt.Items) < 1 {
		fmt.Println("Not enough items")
		return false
	}

	// Checks prices and description of each item
	pricePattern := "\\d+\\.\\d{2}$"
	descPattern := "^[\\w\\s\\-]+$"
	rePrice := regexp.MustCompile(pricePattern)
	reDesc := regexp.MustCompile(descPattern)
	for _, item := range receipt.Items {
		// Price validity
		valid := rePrice.MatchString(item.Price)
		if !valid {
			fmt.Println("Issue with price format")
			return false
		}
		// Description validity
		valid = reDesc.MatchString(item.ShortDescription)
		if !valid {
			fmt.Println("Issue with description format")
			return false
		}
	}
	// All items valid
	return true
}

// Returns unique ID
func GenerateID() string {
	id := uuid.New()
	return id.String()
}

// Handles routing, listens on localhost:8000
func main() {
	router := mux.NewRouter()

	// GET method to get points given a valid receipt ID
	router.HandleFunc("/receipts/process", CreateReceipt).Methods("POST")

	// POST method to create receipt given valid JSON
	router.HandleFunc("/receipts/{id}/points", GetReceiptByID).Methods("GET")

	http.ListenAndServe(":8000", router)

}
