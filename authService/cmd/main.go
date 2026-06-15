package main

import (
	"fmt"
	"net/http"
)

var validUserTokens = map[string]bool{
	"user_key_alpha": true,
	"user_key_beta":  true,
	"user_key_gamma": true,
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	userToken := r.Header.Get("X-User-Token")

	if valid, exists := validUserTokens[userToken]; exists && valid {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Unauthorized User Key", http.StatusUnauthorized)
}

func main() {
	http.HandleFunc("/validate", validateHandler)

	fmt.Println("Auth validation service running on port 5000...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
