# Expense Calculator

A simple web application that helps users track and analyze their monthly expenses.

## Features

- Add new expenses to a list
- Calculate the total amount of expenses
- Calculate the average daily expense (total/30)
- Display the top 3 largest expenses

## Requirements

- Go 1.15 or higher
- Node.js 14+ (for running tests)

## How to Run

1. Clone the repository
2. Navigate to the expense-calculator directory:
   ```bash
   cd app/expense-calculator
   ```
3. Run the Go server:
   ```bash
   go run main.go
   ```
4. Open your browser and go to http://localhost:8080

## Usage

1. Enter expense name and amount in the form
2. Click "Add Expense" to add it to your list
3. View your expense list in the table
4. Delete any expense by clicking the "Delete" button
5. See the summary statistics at the bottom of the page:
   - Total amount of expenses
   - Average daily expense (based on 30 days per month)
   - Top 3 largest expenses

## Testing

The application includes three types of tests:

### Go Backend Tests

To run the Go tests:

```bash
cd app/expense-calculator
go test -v
```

### JavaScript Unit Tests (Jest)

To run the JavaScript unit tests:

```bash
npm install  # Install test dependencies
npm test     # Run Jest tests
```

### End-to-End Tests (Cypress)

To run the Cypress end-to-end tests:

```bash
npm install           # Install test dependencies
npm run start         # Start the application in one terminal
npm run test:e2e      # Run Cypress tests in another terminal
```

Or to open Cypress for interactive testing:

```bash
npm run cypress:open
```

### Browser Tests

For quick manual testing of JavaScript functions, you can also open:

```
http://localhost:8080/tests.html
```

## Implementation Details

- Frontend: HTML, CSS, JavaScript
- Backend: Go (to serve static files)
- No database (in-memory storage only)
- Tests: Go tests, Jest for JavaScript, Cypress for E2E testing 