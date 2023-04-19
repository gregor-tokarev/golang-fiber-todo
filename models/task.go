package models

type Task struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	Text      string `json:"text"`
	DueDate   int64  `json:"dueDate"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
	Status    string `json:"status"`
}

type UpdateTask struct {
	Text    string `json:"text" validate:"required"`
	DueDate int64  `json:"due_date" validate:"required"`
	Status  string `json:"status" validate:"required"`
}
