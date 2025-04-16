package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define the port the server will listen on
	port := 8080

	// Set up static file server
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Start the server
	fmt.Printf("Expense Calculator server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
