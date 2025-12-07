package models

type Assignment struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Subject   string `json:"subject"`
	DueDate   string `json:"due_date"` // YYYY-MM-DD
	Completed bool   `json:"completed"`
}
