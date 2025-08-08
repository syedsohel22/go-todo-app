package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// A struct to represent a single to-do item
type Todo struct {
	ID    int
	Title string
	Done  bool
}

var db *sql.DB

func main() {
	// Connect to the database
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/go-db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")

	// Define our routes and handlers
	http.HandleFunc("/", todoHandler)
	http.HandleFunc("/add", addTodoHandler)
	http.HandleFunc("/toggle", toggleTodoHandler)
	http.HandleFunc("/delete", deleteTodoHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// ... more handler functions will go here

// Handles the main page, displaying all to-do items
func todoHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    rows, err := db.Query("SELECT id, title, done FROM todos ORDER BY id DESC")
    if err != nil {
        http.Error(w, "Failed to query todos", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var todos []Todo
    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done); err != nil {
            http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
            return
        }
        todos = append(todos, todo)
    }

    // Pass the todos to our template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Failed to parse template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, todos)
}

// Handles adding a new to-do item
func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO todos (title) VALUES (?)", title)
	if err != nil {
		http.Error(w, "Failed to add todo", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handles toggling a to-do item's status (done/not done)
func toggleTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    id := r.FormValue("id")
    _, err := db.Exec("UPDATE todos SET done = !done WHERE id = ?", id)
    if err != nil {
        http.Error(w, "Failed to toggle todo", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handles deleting a to-do item
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    id := r.FormValue("id")
    _, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
    if err != nil {
        http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}