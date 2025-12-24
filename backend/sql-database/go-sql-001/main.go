package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersTable struct {
	Users []User
}

func selectUsersQuery() ([]User, error) {
	users := []User{}
	noUsers := errors.New("no users to select")

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		fmt.Println("Error querying users: "+err.Error(), http.StatusInternalServerError)

	}
	defer rows.Close()

	for rows.Next() {
		var u User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			fmt.Println("Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return nil, noUsers
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating: "+err.Error(), http.StatusInternalServerError)
		return nil, noUsers
	}

	return users, nil
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var err error

	/* Frontend */
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing HTML: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err = tmpl.Execute(w, ""); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not alowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := selectUsersQuery()
	if err != nil {
		http.Error(w, "Error getting users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var uTable UsersTable
	uTable.Users = users

	/* Frontend */
	tmpl, err := template.ParseFiles("templates/users.html")
	if err != nil {
		http.Error(w, "Error parsing HTML: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err = tmpl.Execute(w, uTable); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
		return
	}

	var uTable UsersTable

	users, err := selectUsersQuery()
	if err != nil {
		http.Error(w, "Error getting users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID: must be a number : "+err.Error(), http.StatusBadRequest)
			return
		}

		for i := range users {
			if users[i].ID == id {
				uTable.Users = append(uTable.Users, users[i])
			}
		}
	}

	/* Frontend */
	tmpl, err := template.ParseFiles("templates/users.html")
	if err != nil {
		http.Error(w, "Error parsing HTML: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err = tmpl.Execute(w, uTable); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

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

	users, err := selectUsersQuery()
	if err != nil {
		http.Error(w, "Error getting users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var last_id int
	new_id := 0
	for i := range users {
		for j := 0; j < i; j++ {
			last_id = users[i].ID

			if last_id >= new_id {
				new_id = last_id + 1
			}
		}
	}

	sqlInsert := `
		INSERT INTO users (id, name, email)
		VALUES ($1, $2, $3)`

	row := db.QueryRow(sqlInsert, new_id, u.Name, u.Email)
	if row.Err() != nil {
		http.Error(w, "Error saving database: "+row.Err().Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Created new user with ID: %d\n", new_id)
	w.WriteHeader(http.StatusCreated)
}

func setupDatabase() {
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
}

func setupServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", indexPage)
	mux.HandleFunc("GET /users/", getUsersHandler)
	mux.HandleFunc("GET /users/{id}/", getUserById)
	mux.HandleFunc("POST /users/register/", userRegistration)

	fmt.Println("Running on http://localhost:8080/...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Error starting server", err)
	}
}

func main() {
	setupDatabase()
	defer db.Close()

	// HTTP Server
	setupServer()
}
