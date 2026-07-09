package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	CreatedAt   string `json:"createdAt"`
}

type PageData struct {
	Tasks []Task
	Task  Task
}
type ErrorPage struct {
	Error string
}

type SuccessMessage struct {
	Message string
}

var ValidStatus = map[string]bool{
	"Pending":     true,
	"In Progress": true,
	"Completed":   true,
	"Cancelled":   true,
}

var ValidPriority = map[string]bool{
	"High":   true,
	"Medium": true,
	"Low":    true,
}
