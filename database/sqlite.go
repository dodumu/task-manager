package database

import (
	"database/sql"
	"fmt"
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

func CreateTask(task models.Task) (models.Task, error) {
	query := `INSERT INTO tasks(Title, Description, Status, Priority)
	VALUES(?, ?, ?, ?)`

	result, err := DB.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
	)
	if err != nil {
		return models.Task{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.Task{}, err
	}
	task.ID = int(id)
	return GetTaskByID(int(id))
}

func GetAllTasks() ([]models.Task, error) {
	allTasks := []models.Task{}
	query := ` SELECT * FROM tasks `
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		allTasks = append(allTasks, task)
	}

	return allTasks, nil
}

func GetTaskByID(ID int) (models.Task, error) {
	query := `SELECT ID, Title, Description, Status, Priority, CreatedAt FROM tasks
	WHERE id = ?`
	var task models.Task
	row := DB.QueryRow(query, ID)
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Task{}, err
	}
	if err != nil {
		return models.Task{}, fmt.Errorf("tasks with id not found")
	}
	return task, nil
}

func UpdateTask(task models.Task) (models.Task, error) {
	query := `UPDATE tasks 
	SET
    Title = ?,
    Description = ?,
    Status = ?,
    Priority = ?
	WHERE id = ?`

	result, err := DB.Exec(query, task.Title, task.Description, task.Status, task.Priority, task.ID)
	if err != nil {
		return models.Task{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return models.Task{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.Task{}, err
	}
	if rows == 0 {
		return models.Task{}, fmt.Errorf("task with id: %d not found", task.ID)
	}
	return GetTaskByID(int(id))
}

func DeleteTask(id int) error {
	query := `DELETE FROM tasks
	WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return fmt.Errorf("task with id: %d not found", id)
	}
	return nil
}
 

func GetFilteredTasks(status, priority string) ([]models.Task, error) {
	query := `
	SELECT ID, Title, Description, Status, Priority, CreatedAt
	FROM tasks
	WHERE 1=1
	`

	args := []any{}

	if status != "" {
		query += " AND Status = ?"
		args = append(args, status)
	}

	if priority != "" {
		query += " AND Priority = ?"
		args = append(args, priority)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}