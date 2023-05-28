package model

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Due_date    string `json:"due_date"`
	Done        bool   `json:"done"`
	Assigned    string `json:"assigned"`
	CardID      int    `json:"card_id"`
	TaskColor   string `json:"task_color"`
}
