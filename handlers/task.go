package handlers

import (
	"database/sql"
	"encoding/json"
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
	_, err := database.CreateTask(task)
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
	path := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/edit")
	taskID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:

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
	case http.MethodPost:

		title := r.FormValue("title")
		if title == "" {
			http.Error(w, "title cannot be empty", http.StatusBadRequest)
			return
		}
		description := r.FormValue("description")
		status := r.FormValue("status")
		priority := r.FormValue("priority")

		ok := models.ValidStatus[status]
		if !ok {
			http.Error(w, "invalid status selected", http.StatusBadRequest)
			return
		}
		ok = models.ValidPriority[priority]
		if !ok {
			http.Error(w, "invalid priiority selected", http.StatusBadRequest)
			return
		}
		task := models.Task{
			ID:          taskID,
			Title:       title,
			Status:      status,
			Description: description,
			Priority:    priority,
		}

		_, err = database.UpdateTask(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/tasks/"+path, http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func RemoveTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/delete")
	taskID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = database.DeleteTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/task", http.StatusSeeOther)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch {

	case strings.HasSuffix(r.URL.Path, "/edit"):
		EditTask(w, r)
	case strings.HasSuffix(r.URL.Path, "/delete"):
		RemoveTask(w, r)
	default:
		ShowTask(w, r)
	}
}

func ListTasksAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")


	tasks, err := database.GetFilteredTasks(status, priority)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}
	utils.WriteJSON(w, http.StatusOK, tasks)
}

func RecieveAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if task.Title == "" {
		http.Error(w, "title cannot be empty", http.StatusBadRequest)
		return
	}

	if !models.ValidStatus[task.Status] {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	if !models.ValidPriority[task.Priority] {
		http.Error(w, "invalid priority", http.StatusBadRequest)
		return
	}
	createdTask, err := database.CreateTask(task)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, createdTask)
}

func ShowTaskAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorPage{
			Error: "method not allowed",
		})
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorPage{
			Error: "id not available",
		})
		return
	}
	target, err := database.GetTaskByID(idNum)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, models.ErrorPage{
			Error: "task with id not found",
		})
	}
	utils.WriteJSON(w, http.StatusOK, target)
}

func UpdateTaskAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorPage{
			Error: "method not accepted",
		})
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorPage{
			Error: "bad request detected",
		})
		return
	}
	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorPage{
			Error: "bad request detected",
		})
		return
	}
	task.ID = idNum

	if strings.TrimSpace(task.Title) == "" {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorPage{
			Error: "name cannot be empty",
		})
		return
	}

	ok := models.ValidStatus[task.Status]
	if !ok {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorPage{
			Error: "invalid status",
		})
		return
	}
	ok = models.ValidPriority[task.Priority]
	if !ok {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorPage{
			Error: "invalid priority",
		})
		return
	}
	updatedTask, err := database.UpdateTask(task)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, models.ErrorPage{
			Error: "internal server error",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedTask)
}

func DeleteTaskAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, models.ErrorPage{
			Error: "method not allowed",
		})
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, models.ErrorPage{
			Error: "bad request detected",
		})
		return
	}
	err = database.DeleteTask(idNum)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, models.ErrorPage{
			Error: "task with id not found",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, models.SuccessMessage{
		Message: "Task deleted successfully",
	})
}
