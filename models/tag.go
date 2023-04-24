package models

type Tag struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	OwnerId int    `json:"owner_id"`
	Tasks   []Task `json:"tasks"`
}

type CreateTagReq struct {
	Name string `json:"name" validate:"required"`
}

type NewTagConfig struct {
	Name    string
	OwnerId int
}

func NewTag(config NewTagConfig) *Tag {
	var tag = &Tag{}

	tag.Name = config.Name
	tag.OwnerId = config.OwnerId

	DB.Create(&tag)

	return tag
}

func FindTagById(id int) *Tag {
	var tag *Tag
	DB.First(&tag, id)
	return tag
}

func FindAllTags(ownerId int) []Tag {
	var tags []Tag
	DB.Where("owner_id = ?", ownerId).Find(&tags)
	return tags
}

func (t *Tag) Save() {
	DB.Save(t)
}
