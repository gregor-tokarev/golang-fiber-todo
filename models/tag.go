package models

type Tag struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	OwnerId int    `json:"owner_id"`
	Tasks   []Task `json:"tasks"`
}

type CreateTag struct {
	Name string `json:"name" validate:"required"`
}
