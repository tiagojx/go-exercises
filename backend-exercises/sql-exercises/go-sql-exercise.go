package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func userRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlInsert := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id`

	id := 0
	err := db.QueryRow(sqlInsert, u.Name, u.Email).Scan(&id)
	if err != nil {
		http.Error(w, "Error on saving database "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Created new user with ID: %d\n", id)
	w.WriteHeader(http.StatusCreated)
}

func main() {
	var err error

	// Busca pela senha na variável de ambiente.
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("Error: environment variable 'DB_PASSWORD' was not found!")
	}
	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5432/postgres?sslmode=disable", dbPassword)

	// Configura e verifica a conexão com o banco de dados.
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	} else if err = db.Ping(); err != nil {
		log.Fatal("Can't connect to database", err)
	}
	fmt.Println("Successfully connected to PostgreSQL database!")

	defer db.Close()

	// Certifica que a tabela 'users' existe.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Error creating table in database.")
	}
	fmt.Println("Table 'users': ok!")

	// HTTP Server
	http.HandleFunc("/register", userRegistration)

	fmt.Println("Running on http://localhost:8080/...")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server", err)
	}
}
