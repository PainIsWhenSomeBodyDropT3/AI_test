package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock server to return test product data
func setupMockServer(t *testing.T, statusCode int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	}))
}

func TestFetchProducts(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedCount  int
		expectedError  bool
		expectedStatus int
	}{
		{
			name:           "Valid response",
			statusCode:     http.StatusOK,
			responseBody:   `[{"id":1,"title":"Test Product","price":10.5,"description":"Test","category":"test","image":"test.jpg","rating":{"rate":4.5,"count":10}}]`,
			expectedCount:  1,
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty array",
			statusCode:     http.StatusOK,
			responseBody:   `[]`,
			expectedCount:  0,
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Server error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"error":"Internal server error"}`,
			expectedCount:  0,
			expectedError:  true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid JSON",
			statusCode:     http.StatusOK,
			responseBody:   `{invalid-json`,
			expectedCount:  0,
			expectedError:  true,
			expectedStatus: http.StatusOK,
		},
	}

	// Save original API URL
	originalURL := apiURL

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock server
			server := setupMockServer(t, tc.statusCode, tc.responseBody)
			defer server.Close()

			// Override API URL to use mock server
			apiURL = server.URL

			// Call function under test
			products, statusCode, err := fetchProducts()

			// Reset API URL
			apiURL = originalURL

			// Check error
			if tc.expectedError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Check status code
			if statusCode != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, statusCode)
			}

			// Check product count
			if len(products) != tc.expectedCount {
				t.Errorf("Expected %d products, got %d", tc.expectedCount, len(products))
			}
		})
	}
}

func TestValidateProducts(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		products       []Product
		expectedErrors int
	}{
		{
			name: "Valid products",
			products: []Product{
				{
					ID:          1,
					Title:       "Test Product",
					Price:       10.5,
					Description: "Test",
					Category:    "test",
					Image:       "test.jpg",
					Rating:      Rating{Rate: 4.5, Count: 10},
				},
			},
			expectedErrors: 0,
		},
		{
			name: "Empty title",
			products: []Product{
				{
					ID:          1,
					Title:       "",
					Price:       10.5,
					Description: "Test",
					Category:    "test",
					Image:       "test.jpg",
					Rating:      Rating{Rate: 4.5, Count: 10},
				},
			},
			expectedErrors: 1,
		},
		{
			name: "Negative price",
			products: []Product{
				{
					ID:          1,
					Title:       "Test Product",
					Price:       -10.5,
					Description: "Test",
					Category:    "test",
					Image:       "test.jpg",
					Rating:      Rating{Rate: 4.5, Count: 10},
				},
			},
			expectedErrors: 1,
		},
		{
			name: "Rating too high",
			products: []Product{
				{
					ID:          1,
					Title:       "Test Product",
					Price:       10.5,
					Description: "Test",
					Category:    "test",
					Image:       "test.jpg",
					Rating:      Rating{Rate: 5.5, Count: 10},
				},
			},
			expectedErrors: 1,
		},
		{
			name: "Multiple errors",
			products: []Product{
				{
					ID:          1,
					Title:       "",
					Price:       -10.5,
					Description: "",
					Category:    "test",
					Image:       "test.jpg",
					Rating:      Rating{Rate: 5.5, Count: -5},
				},
			},
			expectedErrors: 5, // Empty title, negative price, rating too high, empty description, negative count
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call function under test
			errors := validateProducts(tc.products)

			// Check error count
			if len(errors) != tc.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tc.expectedErrors, len(errors))
			}
		})
	}
}
