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

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
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

func getAllUsers(w http.ResponseWriter, r *http.Request) {
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

func getUserId(w http.ResponseWriter, r *http.Request) {
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

func deleteUserId(w http.ResponseWriter, r *http.Request) {
	var err error
	var id = 0

	idStr := r.PathValue("id")
	if idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID: must be a number: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	res, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error deleting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Printf("User %d was deleted.\n", id)
	w.WriteHeader(http.StatusNoContent)
}

func updateUserId(w http.ResponseWriter, r *http.Request) {
	var err error
	var id = 0
	var u User

	idStr := r.PathValue("id")
	if idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID: must be a number: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println(u.Name + " " + u.Email)
	if u.Name == "" || u.Email == "" {
		var users []User
		var currentUser User

		users, err = selectUsersQuery()
		if err != nil {
			http.Error(w, "Error getting users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for i := range users {
			if users[i].ID == id {
				currentUser = users[i]
			}
		}

		if u.Name == "" {
			u.Name = currentUser.Name
		} else if u.Email == "" {
			u.Email = currentUser.Email
		}
	}

	putUpdateLine := `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3`

	res, err := db.Exec(putUpdateLine, u.Name, u.Email, id)
	if err != nil {
		http.Error(w, "Error updating: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Printf("User %d was updated.\n", id)
	w.WriteHeader(http.StatusOK)
}

func userRegistration(w http.ResponseWriter, r *http.Request) {
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

	// GET users
	mux.HandleFunc("GET /users/", getAllUsers)
	mux.HandleFunc("GET /users/{id}/", getUserId)

	// DELETE users
	mux.HandleFunc("DELETE /users/{id}/", deleteUserId)

	// POST users
	mux.HandleFunc("POST /users/register/", userRegistration)
	mux.HandleFunc("GET /users/register/", methodNotAllowed)

	// UPDATE users
	mux.HandleFunc("PATCH /users/{id}/", updateUserId)

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
