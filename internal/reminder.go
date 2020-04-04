package internal

type Reminder struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	DueDate int64  `json:"due_date"`
}
