package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

func RenderTemplate(w http.ResponseWriter, page string, data any) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+page,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
}
