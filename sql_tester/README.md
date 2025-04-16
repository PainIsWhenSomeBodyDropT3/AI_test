# SQL Tester

A command-line tool and web interface for analyzing sales data using SQL queries.

## Overview

This application demonstrates how to perform common SQL analytics tasks on sales data. It creates an SQLite database, populates it with sample sales data, and runs predefined SQL queries to answer business questions about the data.

## Features

- Creates and populates an SQLite database with sample sales data
- Executes SQL queries to analyze the data:
  - Calculate total sales for a specific period
  - Find the top-spending customer
  - Calculate average order value
- Verifies query results against expected values
- Provides example SQL queries for further exploration
- Includes an interactive web interface for exploring the data and running custom queries

## Requirements

- Go 1.15 or higher
- SQLite3

## Dependencies

This project uses the following third-party libraries:
- github.com/mattn/go-sqlite3 - SQLite driver for Go

To install dependencies:

```bash
go get github.com/mattn/go-sqlite3
```

## How to Run

1. Clone the repository
2. Navigate to the sql_tester directory:
   ```bash
   cd app/sql_tester
   ```
3. Install dependencies:
   ```bash
   go get github.com/mattn/go-sqlite3
   ```
4. Run in CLI mode:
   ```bash
   go run *.go
   ```
5. Run with web interface:
   ```bash
   go run *.go -web
   ```
6. Open your browser and navigate to http://localhost:8080

## Command-Line Options

- `-db string`: Path to the SQLite database (default "sales.db")
- `-web`: Run in web server mode
- `-port int`: Port for web server mode (default 8080)

Example:
```bash
go run *.go -web -port 9090
```

## Web Interface

The web interface provides the following features:

- View the database schema and sample data
- Run predefined queries with validation against expected results
- Write and execute custom SQL queries
- View query results in a formatted table

## SQL Queries

The application tests the following SQL queries:

### 1. Total Sales for March 2024

```sql
SELECT SUM(amount) FROM orders WHERE order_date BETWEEN '2024-03-01' AND '2024-03-31'
```

Expected result: 27,000

### 2. Top-spending Customer

```sql
SELECT customer, SUM(amount) AS total_spent 
FROM orders 
GROUP BY customer 
ORDER BY total_spent DESC 
LIMIT 1
```

Expected result: Alice (20,000)

### 3. Average Order Value

```sql
SELECT AVG(amount) FROM orders
```

Expected result: 6,000

## Additional Custom Queries

Here are some additional queries you can try in the web interface:

### Orders by Month

```sql
SELECT 
    strftime('%Y-%m', order_date) AS month,
    SUM(amount) AS total_sales
FROM orders
GROUP BY month
ORDER BY month
```

### Customer Order Count

```sql
SELECT 
    customer,
    COUNT(*) AS order_count,
    SUM(amount) AS total_spent,
    AVG(amount) AS avg_order_value
FROM orders
GROUP BY customer
ORDER BY total_spent DESC
```

### Daily Sales for March

```sql
SELECT 
    order_date,
    SUM(amount) AS daily_sales
FROM orders
WHERE order_date BETWEEN '2024-03-01' AND '2024-03-31'
GROUP BY order_date
ORDER BY order_date
```

## Data Schema

The application creates a single table with the following schema:

```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    customer TEXT,
    amount REAL,
    order_date DATE
);
```

Sample data:

```sql
INSERT INTO orders (customer, amount, order_date) VALUES
('Alice', 5000, '2024-03-01'),
('Bob', 8000, '2024-03-05'),
('Alice', 3000, '2024-03-15'),
('Charlie', 7000, '2024-02-20'),
('Alice', 10000, '2024-02-28'),
('Bob', 4000, '2024-02-10'),
('Charlie', 9000, '2024-03-22'),
('Alice', 2000, '2024-03-30');
```

## Online SQL Explorers

You can also explore this data using online SQLite explorers:

- [SQLite Online](https://sqliteonline.com/)
- [SQL Fiddle](http://sqlfiddle.com/)
- [DB Fiddle](https://www.db-fiddle.com/)

Simply copy the CREATE TABLE and INSERT INTO statements from this README to get started. 