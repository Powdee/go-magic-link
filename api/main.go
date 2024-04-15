package main

import (
	"encoding/json"
	"net/http"
)

// emptyListHandler sends an empty list as JSON
func emptyListHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Check if the method is GET
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"method not allowed"}`))
		return
	}

	// Encode an empty slice and send it as response
	json.NewEncoder(w).Encode([]struct{}{})
}

func main() {
	// Create a new HTTP route binding it to our handler function
	http.HandleFunc("/emptylist", emptyListHandler)

	// Start the server on port 8080 and log errors if any
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
