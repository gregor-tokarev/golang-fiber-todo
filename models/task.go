package models

type Task struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	Text      string `json:"text"`
	DueDate   int64  `json:"due_date"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	Status    string `json:"status"`
	OwnerId   int    `json:"owner_id"`
}

type UpdateTask struct {
	Text    string `json:"text" validate:"required"`
	DueDate int64  `json:"due_date" validate:"required"`
	Status  string `json:"status" validate:"required"`
}
