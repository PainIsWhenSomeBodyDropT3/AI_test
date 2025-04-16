package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Product data with intentional defects
var mockProducts = []Product{
	{
		ID:          1,
		Title:       "Good Product",
		Price:       19.99,
		Description: "This is a good product with no defects",
		Category:    "good",
		Image:       "good.jpg",
		Rating:      Rating{Rate: 4.5, Count: 120},
	},
	{
		ID:          2,
		Title:       "", // Empty title (defect)
		Price:       29.99,
		Description: "This product has an empty title",
		Category:    "defective",
		Image:       "empty-title.jpg",
		Rating:      Rating{Rate: 3.5, Count: 80},
	},
	{
		ID:          3,
		Title:       "Negative Price Product",
		Price:       -9.99, // Negative price (defect)
		Description: "This product has a negative price",
		Category:    "defective",
		Image:       "negative-price.jpg",
		Rating:      Rating{Rate: 2.5, Count: 40},
	},
	{
		ID:          4,
		Title:       "High Rating Product",
		Price:       39.99,
		Description: "This product has a rating over 5",
		Category:    "defective",
		Image:       "high-rating.jpg",
		Rating:      Rating{Rate: 5.5, Count: 60}, // Rating > 5 (defect)
	},
	{
		ID:          5,
		Title:       "Multiple Problems",
		Price:       -19.99, // Negative price (defect)
		Description: "",     // Empty description (defect)
		Category:    "defective",
		Image:       "multiple-problems.jpg",
		Rating:      Rating{Rate: 6.0, Count: -10}, // Rating > 5 and negative count (defects)
	},
}

// RunMockServer starts a mock server with defective product data
func RunMockServer(port int) {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode the mock products to JSON
		json.NewEncoder(w).Encode(mockProducts)
	})

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting mock server at http://localhost%s/products\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
