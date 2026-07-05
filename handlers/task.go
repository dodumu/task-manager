package handlers

import (
	"net/http"
	"task-manager/utils"
)

func ShowCreateTask(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "create.html", nil)
}
