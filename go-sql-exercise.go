package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{}

func helloWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Web!")
}

func userRegistration(w http.ResponseWriter, r *http.Request) {
	// Verifica o tipo de Request (Method) que o cliente está 'usando'.
	switch r.Method {
	case "GET":
		// Define o cabeçalho (Header) da Response.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users) // envia os usuários (users[]).
	case "POST":
		var newUser User
		json.NewDecoder(r.Body).Decode(&newUser)

		users = append(users, newUser)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(users)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", helloWeb)
	http.HandleFunc("/register", userRegistration)

	fmt.Println("Running on http://localhost:8080/...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
