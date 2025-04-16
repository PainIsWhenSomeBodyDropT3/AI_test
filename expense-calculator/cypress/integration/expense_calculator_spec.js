/// <reference types="cypress" />

describe('Expense Calculator', () => {
  beforeEach(() => {
    // Visit the app before each test
    cy.visit('/');
  });

  it('should load the expense calculator application', () => {
    // Check that the main elements are present
    cy.get('.container h1').should('contain', 'Expense Calculator');
    cy.get('.expense-form').should('be.visible');
    cy.get('.expense-list').should('be.visible');
    cy.get('.summary').should('be.visible');
  });

  it('should add a new expense', () => {
    // Enter expense details
    cy.get('#expense-name').type('Test Expense');
    cy.get('#expense-amount').type('100');
    
    // Submit the form
    cy.get('#expense-form').submit();
    
    // Check the expense was added to the table
    cy.get('#expense-table tr').should('have.length', 1);
    cy.get('#expense-table tr td:first-child').should('contain', 'Test Expense');
    cy.get('#expense-table tr td:nth-child(2)').should('contain', '$100.00');
    
    // Check the summary was updated
    cy.get('#total-amount').should('contain', '$100.00');
    cy.get('#avg-daily').should('contain', '$3.33'); // $100/30 ≈ $3.33
    cy.get('#top-expenses li:first-child').should('contain', 'Test Expense ($100.00)');
  });

  it('should add multiple expenses and calculate summaries correctly', () => {
    // Add first expense
    cy.get('#expense-name').type('Rent');
    cy.get('#expense-amount').type('1000');
    cy.get('#expense-form').submit();
    
    // Add second expense
    cy.get('#expense-name').type('Groceries');
    cy.get('#expense-amount').type('300');
    cy.get('#expense-form').submit();
    
    // Add third expense
    cy.get('#expense-name').type('Utilities');
    cy.get('#expense-amount').type('150');
    cy.get('#expense-form').submit();
    
    // Add fourth expense (smaller than third)
    cy.get('#expense-name').type('Coffee');
    cy.get('#expense-amount').type('50');
    cy.get('#expense-form').submit();
    
    // Check the table has 4 expenses
    cy.get('#expense-table tr').should('have.length', 4);
    
    // Check the total expense is calculated correctly
    cy.get('#total-amount').should('contain', '$1,500.00');
    
    // Check the average daily expense is calculated correctly
    cy.get('#avg-daily').should('contain', '$50.00'); // $1500/30 = $50
    
    // Check the top 3 expenses are displayed in correct order
    cy.get('#top-expenses li:nth-child(1)').should('contain', 'Rent ($1,000.00)');
    cy.get('#top-expenses li:nth-child(2)').should('contain', 'Groceries ($300.00)');
    cy.get('#top-expenses li:nth-child(3)').should('contain', 'Utilities ($150.00)');
  });

  it('should delete expenses', () => {
    // Add an expense
    cy.get('#expense-name').type('Test Expense');
    cy.get('#expense-amount').type('100');
    cy.get('#expense-form').submit();
    
    // Delete the expense
    cy.get('.delete-btn').click();
    
    // Check the expense was removed from the table
    cy.get('#expense-table tr').should('have.length', 0);
    
    // Check the summary was updated
    cy.get('#total-amount').should('contain', '$0.00');
    cy.get('#avg-daily').should('contain', '$0.00');
    cy.get('#top-expenses').should('contain', 'No expenses yet');
  });
}); 