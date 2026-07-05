package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"task-manager/database"
	"task-manager/models"
	"task-manager/utils"
)

func ShowCreateTask(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "create.html", nil)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}
	title := r.FormValue("title")
	description := r.FormValue("description")
	status := r.FormValue("status")
	priority := r.FormValue("priority")

	if strings.TrimSpace(title) == "" {
		http.Error(w, "Task title cannot be empty", http.StatusBadRequest)
		return
	}

	ok := models.ValidStatus[status]
	if !ok {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}
	ok = models.ValidPriority[priority]
	if !ok {
		http.Error(w, "invalid priority selected", http.StatusBadRequest)
		return
	}
	task := models.Task{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
	}
	err := database.CreateTask(task)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/task", http.StatusSeeOther)
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	allTasks := models.PageData{}
	tasks, err := database.GetAllTasks()

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	allTasks.Tasks = tasks
	utils.RenderTemplate(w, "tasks.html", allTasks)
}

func ShowTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")

	taskId, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := database.GetTaskByID(taskId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	pageTask := models.PageData{}
	pageTask.Task = task
	utils.RenderTemplate(w, "show.html", pageTask)
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/edit")
	taskID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := database.GetTaskByID(taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	pageTask := models.PageData{}
	pageTask.Task = task
	utils.RenderTemplate(w, "edit.html", pageTask)
}
