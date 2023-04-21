package models

type Task struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	Text      string `json:"text"`
	DueDate   int64  `json:"due_date"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	Notes     string `json:"notes"`
	Status    string `json:"status"`
	OwnerId   int    `json:"owner_id"`
	Order     int    `json:"order"`
	TagId     int    `json:"tag_id"`
}

type UpdateTask struct {
	Text    string `json:"text"`
	DueDate int64  `json:"due_date"`
	Notes   string `json:"notes"`
}

type ChangeTaskStatus struct {
	Status string `json:"status" validate:"required"`
}

type ChangeTaskOrder struct {
	Order int `json:"order" validate:"required,number"`
}
