package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// WebServer starts a web server that provides a GUI for SQL testing
func runWebServer(db *sql.DB, port int) {
	// Define a handler for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(w, r, db)
	})

	// Define a handler for executing queries
	http.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
		executeQuery(w, r, db)
	})

	// Define a handler for getting table information
	http.HandleFunc("/api/tables", func(w http.ResponseWriter, r *http.Request) {
		getTableInfo(w, r, db)
	})

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting web server at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// serveHome renders the home page
func serveHome(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Prepare template data
	data := struct {
		Queries []Query
	}{
		Queries: []Query{
			{
				Name:        "Total Sales for March 2024",
				Description: "Calculate the total sales volume for March 2024",
				SQL:         "SELECT SUM(amount) FROM orders WHERE order_date BETWEEN '2024-03-01' AND '2024-03-31'",
				Expected:    27000.0,
			},
			{
				Name:        "Top-spending Customer",
				Description: "Find the customer who spent the most overall",
				SQL: `SELECT customer, SUM(amount) AS total_spent 
FROM orders 
GROUP BY customer 
ORDER BY total_spent DESC 
LIMIT 1`,
				Expected: 20000.0,
			},
			{
				Name:        "Average Order Value",
				Description: "Calculate the average order value for all orders",
				SQL:         "SELECT AVG(amount) FROM orders",
				Expected:    6000.0,
			},
		},
	}

	// Define the HTML template
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SQL Tester - Sales Data Analysis</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        h1 {
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        h2 {
            color: #2c3e50;
            margin-top: 30px;
        }
        pre {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
        }
        .query-container {
            margin-bottom: 20px;
            padding: 15px;
            border: 1px solid #ddd;
            border-radius: 5px;
            background-color: #f9f9f9;
        }
        .query-title {
            font-weight: bold;
            margin-bottom: 10px;
        }
        .query-description {
            margin-bottom: 10px;
            color: #555;
        }
        .query-sql {
            font-family: monospace;
            background-color: #f0f0f0;
            padding: 10px;
            border-radius: 5px;
            margin-bottom: 10px;
            white-space: pre-wrap;
        }
        button {
            background-color: #3498db;
            color: white;
            border: none;
            padding: 8px 15px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        button:hover {
            background-color: #2980b9;
        }
        .results {
            margin-top: 15px;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            background-color: #fff;
            display: none;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 10px;
        }
        th, td {
            padding: 8px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f2f2f2;
        }
        .success {
            color: green;
            font-weight: bold;
        }
        .error {
            color: red;
            font-weight: bold;
        }
        .custom-query {
            margin-top: 30px;
            padding: 20px;
            border: 1px solid #ddd;
            border-radius: 5px;
            background-color: #f9f9f9;
        }
        textarea {
            width: 100%;
            min-height: 100px;
            padding: 10px;
            margin-bottom: 10px;
            font-family: monospace;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        .schema-info {
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <h1>SQL Tester - Sales Data Analysis</h1>
    
    <div class="schema-info">
        <h2>Database Schema</h2>
        <pre>
CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    customer TEXT,
    amount REAL,
    order_date DATE
);

Sample data:
INSERT INTO orders (customer, amount, order_date) VALUES
('Alice', 5000, '2024-03-01'),
('Bob', 8000, '2024-03-05'),
('Alice', 3000, '2024-03-15'),
('Charlie', 7000, '2024-02-20'),
('Alice', 10000, '2024-02-28'),
('Bob', 4000, '2024-02-10'),
('Charlie', 9000, '2024-03-22'),
('Alice', 2000, '2024-03-30');</pre>
    </div>

    <h2>Predefined Queries</h2>
    {{range .Queries}}
    <div class="query-container">
        <div class="query-title">{{.Name}}</div>
        <div class="query-description">{{.Description}}</div>
        <div class="query-sql">{{.SQL}}</div>
        <div>Expected Result: {{printf "%.2f" .Expected}}</div>
        <button onclick="runQuery(this, '{{.SQL}}', {{.Expected}})">Run Query</button>
        <div class="results"></div>
    </div>
    {{end}}

    <div class="custom-query">
        <h2>Custom Query</h2>
        <p>Write your own SQL query below:</p>
        <textarea id="custom-sql" placeholder="SELECT * FROM orders LIMIT 10;"></textarea>
        <button onclick="runCustomQuery()">Run Custom Query</button>
        <div class="results" id="custom-results"></div>
    </div>

    <script>
        function runQuery(button, sql, expected) {
            const resultsDiv = button.nextElementSibling;
            resultsDiv.style.display = 'block';
            resultsDiv.innerHTML = 'Loading...';
            
            fetch('/api/query', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query: sql }),
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    resultsDiv.innerHTML = '<div class="error">Error: ' + data.error + '</div>';
                    return;
                }
                
                let html = '<h3>Results:</h3>';
                
                // Create table for results
                if (data.results && data.results.length > 0) {
                    html += '<table><thead><tr>';
                    
                    // Table headers from the first result's keys
                    for (const key in data.results[0]) {
                        html += '<th>' + key + '</th>';
                    }
                    html += '</tr></thead><tbody>';
                    
                    // Table rows
                    for (const row of data.results) {
                        html += '<tr>';
                        for (const key in row) {
                            html += '<td>' + row[key] + '</td>';
                        }
                        html += '</tr>';
                    }
                    html += '</tbody></table>';
                    
                    // Check if result matches expected value
                    let actualValue = null;
                    if (data.results[0].hasOwnProperty('SUM(amount)')) {
                        actualValue = data.results[0]['SUM(amount)'];
                    } else if (data.results[0].hasOwnProperty('AVG(amount)')) {
                        actualValue = data.results[0]['AVG(amount)'];
                    } else if (data.results[0].hasOwnProperty('total_spent')) {
                        actualValue = data.results[0]['total_spent'];
                    }
                    
                    if (actualValue !== null && Math.abs(actualValue - expected) < 0.01) {
                        html += '<div class="success">✓ Result matches expected value!</div>';
                    } else if (actualValue !== null) {
                        html += '<div class="error">✗ Result does not match expected value. Expected: ' + expected.toFixed(2) + '</div>';
                    }
                } else {
                    html += '<p>No results returned</p>';
                }
                
                resultsDiv.innerHTML = html;
            })
            .catch(error => {
                resultsDiv.innerHTML = '<div class="error">Error: ' + error.message + '</div>';
            });
        }
        
        function runCustomQuery() {
            const sql = document.getElementById('custom-sql').value;
            const resultsDiv = document.getElementById('custom-results');
            resultsDiv.style.display = 'block';
            resultsDiv.innerHTML = 'Loading...';
            
            fetch('/api/query', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query: sql }),
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    resultsDiv.innerHTML = '<div class="error">Error: ' + data.error + '</div>';
                    return;
                }
                
                let html = '<h3>Results:</h3>';
                
                // Create table for results
                if (data.results && data.results.length > 0) {
                    html += '<table><thead><tr>';
                    
                    // Table headers from the first result's keys
                    for (const key in data.results[0]) {
                        html += '<th>' + key + '</th>';
                    }
                    html += '</tr></thead><tbody>';
                    
                    // Table rows
                    for (const row of data.results) {
                        html += '<tr>';
                        for (const key in row) {
                            html += '<td>' + row[key] + '</td>';
                        }
                        html += '</tr>';
                    }
                    html += '</tbody></table>';
                } else {
                    html += '<p>No results returned</p>';
                }
                
                resultsDiv.innerHTML = html;
            })
            .catch(error => {
                resultsDiv.innerHTML = '<div class="error">Error: ' + error.message + '</div>';
            });
        }
    </script>
</body>
</html>
`

	// Parse and execute the template
	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

// executeQuery executes a SQL query and returns the results as JSON
func executeQuery(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var request struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing request: %v", err), http.StatusBadRequest)
		return
	}

	// Execute the query
	rows, err := db.Query(request.Query)
	if err != nil {
		// Return error as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Error getting columns: %v", err),
		})
		return
	}

	// Prepare a slice to hold the results
	var results []map[string]interface{}

	// Prepare a slice to hold values for each row
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Iterate through rows
	for rows.Next() {
		// Scan the row into the slice of pointers
		err := rows.Scan(valuePtrs...)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Error scanning row: %v", err),
			})
			return
		}

		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			// Convert to appropriate type and store in the map
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}

		// Add the row to the results
		results = append(results, row)
	}

	// Check for errors in rows.Next()
	if err := rows.Err(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Error iterating rows: %v", err),
		})
		return
	}

	// Return the results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

// getTableInfo returns information about the tables in the database
func getTableInfo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Query to get table names
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Error getting tables: %v", err),
		})
		return
	}
	defer rows.Close()

	// Collect table names
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Error scanning table name: %v", err),
			})
			return
		}
		tables = append(tables, name)
	}

	// Get schema for each table
	var tableInfo []map[string]interface{}
	for _, tableName := range tables {
		// Skip internal SQLite tables
		if strings.HasPrefix(tableName, "sqlite_") {
			continue
		}

		// Get table schema
		schemaRows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
		if err != nil {
			continue
		}
		defer schemaRows.Close()

		var columns []map[string]interface{}
		for schemaRows.Next() {
			var cid int
			var name, typeName string
			var notNull, pk int
			var dfltValue interface{}

			if err := schemaRows.Scan(&cid, &name, &typeName, &notNull, &dfltValue, &pk); err != nil {
				continue
			}

			columns = append(columns, map[string]interface{}{
				"name":     name,
				"type":     typeName,
				"not_null": notNull == 1,
				"pk":       pk == 1,
				"default":  dfltValue,
			})
		}

		tableInfo = append(tableInfo, map[string]interface{}{
			"name":    tableName,
			"columns": columns,
		})
	}

	// Return the table info as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tables": tableInfo,
	})
}
