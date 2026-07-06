package main

import (
	"fmt"
	"log"
	"net/http"
	"task-manager/database"
	"task-manager/handlers"
)

func main() {
	database.Init()
	defer database.DB.Close()

	query := `CREATE TABLE IF NOT EXISTS tasks (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	Title TEXT NOT NULL, 
	Description TEXT,
	Status TEXT NOT NULL, 
	Priority TEXT NOT NULL,
	CreatedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server is running on port: http://localhost:8085")
	http.HandleFunc("/tasks/new", handlers.ShowCreateTask)
	http.HandleFunc("/tasks", handlers.CreateTask)
	http.HandleFunc("/task", handlers.ListTasks)
	http.HandleFunc("/tasks/", handlers.TaskHandler)

	http.ListenAndServe(":8085", nil)
}
