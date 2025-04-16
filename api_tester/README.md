# API Tester

A command-line tool for validating data from the FakeStore API to detect errors and anomalies.

## Features

- Verifies server response code (expected 200)
- Validates product data attributes:
  - `title` (name) - must not be empty
  - `price` (price) - must not be negative
  - `rating.rate` - must not exceed 5
- Additional validations:
  - Zero price detection
  - Empty description detection
  - Negative rating count detection
- Generates detailed reports in console or JSON format
- Provides formatted tabular output of defects
- Includes a mock server with intentionally defective data for testing

## Requirements

- Go 1.15 or higher

## How to Run

1. Clone the repository
2. Navigate to the api_tester directory:
   ```bash
   cd app/api_tester
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
4. For JSON output, use the `-json` flag:
   ```bash
   go run main.go -json report.json
   ```
5. To test against mock data with defects:
   ```bash
   go run main.go -mock
   ```
6. You can combine options:
   ```bash
   go run main.go -mock -port 9090 -json mock-report.json
   ```

## Testing

Run the unit tests:

```bash
cd app/api_tester
go test -v
```

## Mock Server

The API tester includes a mock server that serves product data with intentional defects for testing purposes. The mock server includes the following defects:

- Product with empty title
- Product with negative price
- Product with rating > 5
- Product with multiple defects (negative price, empty description, rating > 5, negative rating count)

To use the mock server:

```bash
go run main.go -mock
```

This will start the mock server and run the API tester against it, allowing you to see how the validation works with defective data.

## Example Output

```
API Tester - FakeStore API Validation
=====================================

Testing API: http://localhost:8080/products

Test 1: Verify server response code
Status Code: 200
✅ Status code is 200 OK

Test 2: Validate product attributes
Total products: 5
Products with defects: 7

Defective Products:
-----------------
ID  Title                       Field        Issue                  Value
--  -----                       -----        -----                  -----
2   <empty>                     title        Title is empty         
3   Negative Price Product      price        Price is negative      -9.99
4   High Rating Product         rating.rate  Rating rate exceeds 5  5.5
5   Multiple Problems           price        Price is negative      -19.99
5   Multiple Problems           description  Description is empty   
5   Multiple Problems           rating.rate  Rating rate exceeds 5  6
5   Multiple Problems           rating.count Rating count is negative -10
```

## Implementation Details

The API Tester checks the [FakeStore API](https://fakestoreapi.com/products) for data quality issues using the following process:

1. Fetches product data from the API
2. Validates the HTTP response code
3. Validates each product against a set of rules
4. Generates a report of defects found
5. Optionally outputs a detailed JSON report

The tool is designed to detect anomalies in API data that might cause issues in applications consuming the API. 