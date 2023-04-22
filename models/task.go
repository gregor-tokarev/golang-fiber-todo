package models

import "gorm.io/gorm"

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

type UpdateTaskReq struct {
	Text    string `json:"text"`
	DueDate int64  `json:"due_date"`
	Notes   string `json:"notes"`
}

type ChangeTaskStatusReq struct {
	Status string `json:"status" validate:"required"`
}

type ChangeTaskOrderReq struct {
	Order int `json:"order" validate:"required,number"`
}

type CreateTaskConfig struct {
	OwnerId int `json:"owner_id"`
}

func NewTask(config CreateTaskConfig) *Task {
	var task Task
	task.Status = "todo"
	task.Text = ""

	maxOrder := getMaxTaskOrder(config.OwnerId)
	task.Order = maxOrder + 1

	task.OwnerId = config.OwnerId

	DB.Debug().Create(&task)

	return &task
}

func FindTaskById(id int) *Task {
	var task *Task
	DB.First(&task, id)
	return task
}

type FindTasksConfig struct {
	OwnerId int
	TagId   string
	Skip    int
	Take    int
}

func FindTasks(config FindTasksConfig) []Task {
	var tasks []Task
	dbReq := DB.Where("owner_id = ?", config.OwnerId).Offset(config.Skip).Limit(config.Take)
	if config.TagId != "" {
		dbReq = dbReq.Where("tag_id = ?", config.TagId)
	}

	dbReq.Find(&tasks)

	return tasks
}

func (t *Task) Delete() {
	DB.Delete(&Task{}, t.Id)
}

func (t *Task) Save() {
	DB.Save(&t)
}

func getMaxTaskOrder(userId int) int {
	var maxOrder int
	DB.Raw("SELECT max(tasks.order) FROM tasks WHERE owner_id = ?", userId).Scan(&maxOrder)
	return maxOrder
}

func (t *Task) ChangeTaskOrder(order int) {
	if t.Order > order {
		DB.
			Model(&Task{}).
			Where("\"owner_id\" = ? AND \"order\" >= ? AND \"order\" <= ?", t.OwnerId, order, t.Order).
			Update("\"order\"", gorm.Expr("\"order\" + 1"))
	} else if t.Order < order {
		DB.
			Model(&Task{}).
			Where("\"owner_id\" = ? AND \"order\" <= ? AND \"order\" > ?", t.OwnerId, order, t.Order).
			Update("\"order\"", gorm.Expr("\"order\" - 1"))
	}

	t.Order = order
	t.Save()
}

func (t *Task) ClearTag() {
	DB.Model(&Task{}).Update("tag_id", gorm.Expr("NULL"))
}
