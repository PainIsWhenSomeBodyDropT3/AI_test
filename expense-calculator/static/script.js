// Initialize expenses array
let expenses = [];

// DOM Elements
const expenseForm = document.getElementById('expense-form');
const expenseNameInput = document.getElementById('expense-name');
const expenseAmountInput = document.getElementById('expense-amount');
const expenseTable = document.getElementById('expense-table');
const totalAmountElement = document.getElementById('total-amount');
const avgDailyElement = document.getElementById('avg-daily');
const topExpensesElement = document.getElementById('top-expenses');

// Days in month (default to 30)
const daysInMonth = 30;

// Add event listener to the form
expenseForm.addEventListener('submit', function(e) {
    e.preventDefault();
    
    // Get values from inputs
    const name = expenseNameInput.value.trim();
    const amountStr = expenseAmountInput.value.trim();
    
    // Validate inputs
    if (name === '' || amountStr === '') {
        alert('Please enter both name and amount');
        return;
    }
    
    // Convert amount to number and validate
    const amount = parseFloat(amountStr);
    if (isNaN(amount) || amount <= 0) {
        alert('Please enter a valid amount greater than 0');
        return;
    }
    
    // Add expense to the array
    addExpense(name, amount);
    
    // Clear form inputs
    expenseNameInput.value = '';
    expenseAmountInput.value = '';
    expenseNameInput.focus();
});

// Function to add an expense
function addExpense(name, amount) {
    // Create expense object
    const expense = {
        id: Date.now(), // Using timestamp as a unique id
        name: name,
        amount: amount
    };
    
    // Add to expenses array
    expenses.push(expense);
    
    // Update UI
    displayExpenses();
    updateSummary();
}

// Function to remove an expense
function removeExpense(id) {
    expenses = expenses.filter(expense => expense.id !== id);
    
    // Update UI
    displayExpenses();
    updateSummary();
}

// Function to display expenses in the table
function displayExpenses() {
    // Clear the table first
    expenseTable.innerHTML = '';
    
    // Add each expense to the table
    expenses.forEach(expense => {
        const row = document.createElement('tr');
        
        row.innerHTML = `
            <td>${expense.name}</td>
            <td>$${expense.amount.toFixed(2)}</td>
            <td><button class="delete-btn" data-id="${expense.id}">Delete</button></td>
        `;
        
        expenseTable.appendChild(row);
    });
    
    // Add event listeners to delete buttons
    document.querySelectorAll('.delete-btn').forEach(button => {
        button.addEventListener('click', function() {
            const id = parseInt(this.getAttribute('data-id'));
            removeExpense(id);
        });
    });
}

// Function to update the summary section
function updateSummary() {
    // Calculate total amount
    const totalAmount = expenses.reduce((total, expense) => total + expense.amount, 0);
    
    // Calculate average daily expense
    const avgDaily = totalAmount / daysInMonth;
    
    // Get top 3 largest expenses
    const topExpenses = [...expenses]
        .sort((a, b) => b.amount - a.amount)
        .slice(0, 3);
    
    // Update UI elements
    totalAmountElement.textContent = `$${totalAmount.toFixed(2)}`;
    avgDailyElement.textContent = `$${avgDaily.toFixed(2)}`;
    
    // Update top expenses list
    if (topExpenses.length === 0) {
        topExpensesElement.innerHTML = '<li>No expenses yet</li>';
    } else {
        topExpensesElement.innerHTML = topExpenses
            .map(expense => `<li>${expense.name} ($${expense.amount.toFixed(2)})</li>`)
            .join('');
    }
}

// Initialize the application
displayExpenses();
updateSummary(); 