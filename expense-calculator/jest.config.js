module.exports = {
  testEnvironment: 'jsdom',
  testMatch: ['**/static/**/*.test.js'],
  collectCoverage: true,
  coverageDirectory: 'coverage',
  coverageReporters: ['text', 'lcov'],
  verbose: true,
}; 