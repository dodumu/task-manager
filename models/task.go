package models

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Priority    string
	CreatedAt   string
}

type PageData struct {
	Tasks []Task
	Task  Task
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
