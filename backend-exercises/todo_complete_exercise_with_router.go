package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const SERVER_PORT string = ":8080"

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{ID: 1, Title: "Aprender Go", Done: false},
	{ID: 2, Title: "Aprender 'net/http'", Done: true},
	{ID: 3, Title: "Aprender REST e go-chi", Done: false},
}

var currentID = 3

func helloWeb(w http.ResponseWriter, h *http.Request) { fmt.Fprintf(w, "Hello, Web") }

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	case "POST":
		var newTask Task
		json.NewDecoder(r.Body).Decode(&newTask)

		currentID++
		newTask.ID = currentID
		tasks = append(tasks, newTask)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func findTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		var itemFound bool
		var itemById Task

		for _, i := range tasks {
			if i.ID == id {
				itemById = i
				itemFound = true
			}
		}

		if itemFound {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(itemById)
		} else {
			http.Error(w, "Item not found", http.StatusNotFound)
		}
	case "PUT":
		var overrideTask Task
		json.NewDecoder(r.Body).Decode(&overrideTask)

		tasks[id] = overrideTask

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tasks[id])
	case "DELETE":
		tasks := append(tasks[:id], tasks[id+1:]...)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	r := chi.NewRouter()

	r.HandleFunc("/", helloWeb)
	r.Method("GET", "/tasks", http.HandlerFunc(tasksHandler))
	r.Method("POST", "/tasks", http.HandlerFunc(tasksHandler))
	r.Method("GET", "/tasks/{id}", http.HandlerFunc(findTaskByIdHandler))

	fmt.Printf("Starting server on port %s...\n", SERVER_PORT)
	http.ListenAndServe(SERVER_PORT, r)
}
