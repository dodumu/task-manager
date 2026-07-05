package database

import (
	"database/sql"
	"log"
	"task-manager/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error

	DB, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func CreateTask(task models.Task) error {
	query := `INSERT INTO tasks(Title, Description, Status, Priority)
	VALUES(?, ?, ?, ?)`

	_, err := DB.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
	)
	return err
}
