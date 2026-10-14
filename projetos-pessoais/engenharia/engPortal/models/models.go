package models

type ProjectResponse struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	LastCommitDate string `json:"last_commit_date"`
	Branch         string `json:"branch"`
	StatusCI       string `json:"status_ci"` // Ex: Success, Running, Failed
	WebURL         string `json:"web_url"`
}

type KanbanTask struct {
	Title       string `json:"title"`
	CreatedBy   string `json:"createdBy"`
	DueDateTime string `json:"dueDateTime"`
	Status      string `json:"status"`
	BucketID    string `json:"-"` // Internal use for grouping
}

type KanbanBucket struct {
	ID    string       `json:"id"`
	Name  string       `json:"name"`
	Tasks []KanbanTask `json:"tasks"`
}

type KanbanBoard struct {
	PlanID  string         `json:"planId"`
	Title   string         `json:"title"`
	Buckets []KanbanBucket `json:"buckets"`
}
