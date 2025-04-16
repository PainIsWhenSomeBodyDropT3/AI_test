package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Query represents a SQL query with its expected result
type Query struct {
	Name        string
	Description string
	SQL         string
	Expected    float64
}

func main() {
	// Parse command line flags
	dbPath := flag.String("db", "sales.db", "Path to the SQLite database")
	webMode := flag.Bool("web", false, "Run in web server mode")
	webPort := flag.Int("port", 8080, "Port for web server mode")
	flag.Parse()

	fmt.Println("SQL Tester - Sales Data Analysis")
	fmt.Println("================================")
	fmt.Println()

	// Create or open the database
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create and populate the database
	if err := setupDatabase(db); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	fmt.Println("✅ Database created and populated with sample data")
	fmt.Println()

	// Check if web mode is enabled
	if *webMode {
		fmt.Printf("Starting web interface on port %d...\n", *webPort)
		runWebServer(db, *webPort)
		return
	}

	// Define the queries to test
	queries := []Query{
		{
			Name:        "Total Sales for March 2024",
			Description: "Calculate the total sales volume for March 2024",
			SQL:         "SELECT SUM(amount) FROM orders WHERE order_date BETWEEN '2024-03-01' AND '2024-03-31'",
			Expected:    27000.0,
		},
		{
			Name:        "Top-spending Customer",
			Description: "Find the customer who spent the most overall",
			SQL: `
				SELECT customer, SUM(amount) AS total_spent 
				FROM orders 
				GROUP BY customer 
				ORDER BY total_spent DESC 
				LIMIT 1
			`,
			Expected: 20000.0,
		},
		{
			Name:        "Average Order Value",
			Description: "Calculate the average order value for all orders",
			SQL:         "SELECT AVG(amount) FROM orders",
			Expected:    6000.0,
		},
	}

	// Run the queries and check results
	for _, q := range queries {
		fmt.Printf("Test: %s\n", q.Name)
		fmt.Printf("Description: %s\n", q.Description)
		fmt.Printf("SQL Query: %s\n", q.SQL)

		var result float64
		var customer string

		// For the top-spending customer query, we need to handle multiple columns
		if q.Name == "Top-spending Customer" {
			err = db.QueryRow(q.SQL).Scan(&customer, &result)
		} else {
			err = db.QueryRow(q.SQL).Scan(&result)
		}

		if err != nil {
			fmt.Printf("❌ Error executing query: %v\n", err)
			continue
		}

		if q.Name == "Top-spending Customer" {
			fmt.Printf("Result: %s (%.2f)\n", customer, result)
			if result == q.Expected && customer == "Alice" {
				fmt.Println("✅ Result matches expected value")
			} else {
				fmt.Printf("❌ Result does not match expected value. Expected: Alice (%.2f)\n", q.Expected)
			}
		} else {
			fmt.Printf("Result: %.2f\n", result)
			if result == q.Expected {
				fmt.Println("✅ Result matches expected value")
			} else {
				fmt.Printf("❌ Result does not match expected value. Expected: %.2f\n", q.Expected)
			}
		}
		fmt.Println()
	}

	// Output a prompt to view the data
	fmt.Println("To manually inspect the data, you can run the app with the web interface:")
	fmt.Println("- go run *.go -web")
	fmt.Println("Or run these commands in SQLite:")
	fmt.Println("- To see all orders: SELECT * FROM orders;")
	fmt.Println("- To see customer totals: SELECT customer, SUM(amount) FROM orders GROUP BY customer;")
	fmt.Println()

	fmt.Println("All tests completed.")
}

// setupDatabase creates and populates the orders table
func setupDatabase(db *sql.DB) error {
	// Drop the table if it exists
	_, err := db.Exec("DROP TABLE IF EXISTS orders")
	if err != nil {
		return fmt.Errorf("failed to drop table: %w", err)
	}

	// Create the orders table
	_, err = db.Exec(`
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY,
			customer TEXT,
			amount REAL,
			order_date DATE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Insert sample data
	_, err = db.Exec(`
		INSERT INTO orders (customer, amount, order_date) VALUES
		('Alice', 5000, '2024-03-01'),
		('Bob', 8000, '2024-03-05'),
		('Alice', 3000, '2024-03-15'),
		('Charlie', 7000, '2024-02-20'),
		('Alice', 10000, '2024-02-28'),
		('Bob', 4000, '2024-02-10'),
		('Charlie', 9000, '2024-03-22'),
		('Alice', 2000, '2024-03-30')
	`)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}
