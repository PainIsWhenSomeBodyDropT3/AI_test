package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// API endpoint URL (variable for testing)
var apiURL = "https://fakestoreapi.com/products"

// Product represents a product from the FakeStore API
type Product struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Rating      Rating  `json:"rating"`
}

// Rating represents the rating information for a product
type Rating struct {
	Rate  float64 `json:"rate"`
	Count int     `json:"count"`
}

// ValidationError represents an error found during validation
type ValidationError struct {
	ProductID   int         `json:"product_id"`
	Title       string      `json:"title"`
	Field       string      `json:"field"`
	Message     string      `json:"message"`
	ActualValue interface{} `json:"actual_value"`
}

// TestReport represents the overall test results
type TestReport struct {
	Timestamp       string            `json:"timestamp"`
	URL             string            `json:"url"`
	StatusCode      int               `json:"status_code"`
	StatusCodeValid bool              `json:"status_code_valid"`
	TotalProducts   int               `json:"total_products"`
	DefectCount     int               `json:"defect_count"`
	Defects         []ValidationError `json:"defects"`
}

func main() {
	// Parse command line flags
	jsonOutput := flag.String("json", "", "Output JSON report to specified file")
	mockServer := flag.Bool("mock", false, "Run with mock server containing defective data")
	mockPort := flag.Int("port", 8080, "Port for mock server")
	flag.Parse()

	// Run mock server if requested
	if *mockServer {
		// Update URL to point to local mock server
		apiURL = fmt.Sprintf("http://localhost:%d/products", *mockPort)

		// Launch mock server in a goroutine
		go RunMockServer(*mockPort)

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("API Tester - FakeStore API Validation")
	fmt.Println("=====================================")
	fmt.Println()

	// Display which API we're testing
	fmt.Printf("Testing API: %s\n\n", apiURL)

	// Initialize test report
	report := TestReport{
		Timestamp: time.Now().Format(time.RFC3339),
		URL:       apiURL,
	}

	// Fetch data from API
	products, statusCode, err := fetchProducts()
	if err != nil {
		fmt.Printf("Error fetching products: %v\n", err)
		os.Exit(1)
	}

	// Test 1: Verify server response code
	fmt.Println("Test 1: Verify server response code")
	fmt.Printf("Status Code: %d\n", statusCode)
	report.StatusCode = statusCode
	report.StatusCodeValid = (statusCode == http.StatusOK)

	if report.StatusCodeValid {
		fmt.Println("✅ Status code is 200 OK")
	} else {
		fmt.Printf("❌ Expected status code 200, got %d\n", statusCode)
	}
	fmt.Println()

	// Validate products and collect errors
	fmt.Println("Test 2: Validate product attributes")
	validationErrors := validateProducts(products)

	// Update report
	report.TotalProducts = len(products)
	report.DefectCount = len(validationErrors)
	report.Defects = validationErrors

	// Display validation results
	fmt.Printf("Total products: %d\n", report.TotalProducts)
	fmt.Printf("Products with defects: %d\n", report.DefectCount)
	fmt.Println()

	// Display the list of defects
	if report.DefectCount > 0 {
		fmt.Println("Defective Products:")
		fmt.Println("-----------------")
		printValidationErrors(validationErrors)
	} else {
		fmt.Println("✅ No defects found in any products")
	}

	// Output JSON report if requested
	if *jsonOutput != "" {
		generateJSONReport(*jsonOutput, report)
	}

	// If running mock server, don't exit immediately
	if *mockServer {
		fmt.Println("\nMock server is running. Press Ctrl+C to exit.")
		// Block to keep the server running
		select {}
	}
}

// fetchProducts retrieves products from the API
func fetchProducts() ([]Product, int, error) {
	// Make HTTP request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON
	var products []Product
	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return products, resp.StatusCode, nil
}

// validateProducts checks all products for defects
func validateProducts(products []Product) []ValidationError {
	var errors []ValidationError

	for _, product := range products {
		// Check for empty title
		if product.Title == "" {
			errors = append(errors, ValidationError{
				ProductID:   product.ID,
				Title:       product.Title,
				Field:       "title",
				Message:     "Title is empty",
				ActualValue: product.Title,
			})
		} else if strings.TrimSpace(product.Title) == "" {
			errors = append(errors, ValidationError{
				ProductID:   product.ID,
				Title:       product.Title,
				Field:       "title",
				Message:     "Title contains only whitespace",
				ActualValue: product.Title,
			})
		}

		// Check for negative price
		if product.Price < 0 {
			errors = append(errors, ValidationError{
				ProductID:   product.ID,
				Title:       product.Title,
				Field:       "price",
				Message:     "Price is negative",
				ActualValue: product.Price,
			})
		}

		// Check for rating rate > 5
		if product.Rating.Rate > 5 {
			errors = append(errors, ValidationError{
				ProductID:   product.ID,
				Title:       product.Title,
				Field:       "rating.rate",
				Message:     "Rating rate exceeds 5",
				ActualValue: product.Rating.Rate,
			})
		}

		// Add extra validation: check for zero price
		if product.Price == 0 {
			errors = append(errors, ValidationError{
				ProductID:   product.ID,
				Title:       product.Title,
				Field:       "price",
				Message:     "Price is zero",
				ActualValue: product.Price,
			})
		}

		// Add extra validation: check for negative rating count
		if product.Rating.Count < 0 {
			errors = append(errors, ValidationError{
				ProductID:   product.ID,
				Title:       product.Title,
				Field:       "rating.count",
				Message:     "Rating count is negative",
				ActualValue: product.Rating.Count,
			})
		}

		// Add extra validation: check for empty description
		if product.Description == "" {
			errors = append(errors, ValidationError{
				ProductID:   product.ID,
				Title:       product.Title,
				Field:       "description",
				Message:     "Description is empty",
				ActualValue: product.Description,
			})
		}
	}

	return errors
}

// printValidationErrors displays validation errors in a formatted table
func printValidationErrors(errors []ValidationError) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTitle\tField\tIssue\tValue")
	fmt.Fprintln(w, "--\t-----\t-----\t-----\t-----")

	for _, err := range errors {
		title := err.Title
		if title == "" {
			title = "<empty>"
		} else if len(title) > 30 {
			title = title[:27] + "..."
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%v\n",
			err.ProductID,
			title,
			err.Field,
			err.Message,
			err.ActualValue,
		)
	}
	w.Flush()
}

// generateJSONReport outputs the test report as JSON to a file
func generateJSONReport(filename string, report TestReport) {
	// Marshal report to JSON
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON report: %v\n", err)
		return
	}

	// Write to file
	err = os.WriteFile(filename, reportJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON report to file %s: %v\n", filename, err)
		return
	}

	fmt.Printf("\nJSON report written to %s\n", filename)
}
