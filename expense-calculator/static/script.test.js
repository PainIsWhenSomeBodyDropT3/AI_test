/**
 * @jest-environment jsdom
 */

// Mock the DOM elements
document.body.innerHTML = `
<form id="expense-form"></form>
<input id="expense-name">
<input id="expense-amount">
<div id="expense-table"></div>
<span id="total-amount"></span>
<span id="avg-daily"></span>
<ul id="top-expenses"></ul>
`;

// Import the script (since we're mocking browser environment, we need to evaluate the script)
const fs = require('fs');
const path = require('path');
const scriptContent = fs.readFileSync(path.resolve(__dirname, 'script.js'), 'utf8');

// Define global mocks before evaluating the script
global.expenses = [];
global.daysInMonth = 30;

// Mock DOM functions
document.getElementById = jest.fn(id => {
    return {
        addEventListener: jest.fn(),
        focus: jest.fn(),
        innerHTML: '',
        value: '',
        appendChild: jest.fn(),
        querySelectorAll: jest.fn().mockReturnValue([{
            addEventListener: jest.fn(),
            getAttribute: jest.fn().mockReturnValue('1')
        }]),
        textContent: ''
    };
});

// Evaluate the script
eval(scriptContent);

// Make functions globally available for testing
global.addExpense = addExpense;
global.removeExpense = removeExpense;
global.displayExpenses = jest.fn();
global.updateSummary = jest.fn();

describe('Expense Calculator', () => {
    beforeEach(() => {
        // Reset expenses before each test
        global.expenses = [];
        // Reset mocks
        displayExpenses.mockClear();
        updateSummary.mockClear();
    });

    test('addExpense should add an expense to the array', () => {
        addExpense('Test Expense', 100);
        
        expect(expenses.length).toBe(1);
        expect(expenses[0].name).toBe('Test Expense');
        expect(expenses[0].amount).toBe(100);
        expect(displayExpenses).toHaveBeenCalled();
        expect(updateSummary).toHaveBeenCalled();
    });

    test('removeExpense should remove an expense from the array', () => {
        // Add two expenses
        addExpense('Expense 1', 100);
        addExpense('Expense 2', 200);
        
        // Reset mock calls from addExpense
        displayExpenses.mockClear();
        updateSummary.mockClear();
        
        // Get ID of the first expense
        const idToRemove = expenses[0].id;
        
        // Remove the first expense
        removeExpense(idToRemove);
        
        expect(expenses.length).toBe(1);
        expect(expenses[0].name).toBe('Expense 2');
        expect(displayExpenses).toHaveBeenCalled();
        expect(updateSummary).toHaveBeenCalled();
    });

    test('Total calculation should work correctly', () => {
        addExpense('Expense 1', 100);
        addExpense('Expense 2', 200);
        addExpense('Expense 3', 300);
        
        const totalAmount = expenses.reduce((total, expense) => total + expense.amount, 0);
        
        expect(totalAmount).toBe(600);
    });

    test('Average daily calculation should work correctly', () => {
        addExpense('Expense 1', 300);
        
        const totalAmount = expenses.reduce((total, expense) => total + expense.amount, 0);
        const avgDaily = totalAmount / daysInMonth;
        
        expect(avgDaily).toBe(10);
    });

    test('Top expenses should be calculated correctly', () => {
        addExpense('Small', 50);
        addExpense('Medium', 100);
        addExpense('Large', 200);
        addExpense('Extra Large', 300);
        
        const topExpenses = [...expenses]
            .sort((a, b) => b.amount - a.amount)
            .slice(0, 3);
        
        expect(topExpenses.length).toBe(3);
        expect(topExpenses[0].amount).toBe(300);
        expect(topExpenses[1].amount).toBe(200);
        expect(topExpenses[2].amount).toBe(100);
    });
}); 