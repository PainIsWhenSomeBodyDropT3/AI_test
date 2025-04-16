// ***********************************************************
// This is a support file that will be loaded automatically.
// You can change its location using cypress.json
// ***********************************************************

// Import commands.js using ES2015 syntax:
// import './commands'

// Alternatively you can use CommonJS syntax:
// require('./commands')

// Log when tests start/end
beforeEach(() => {
  cy.log('Test starting...');
});

afterEach(() => {
  cy.log('Test complete.');
}); 