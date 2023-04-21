package models

type User struct {
	Id           int    `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken"`
	Tasks        []Task `gorm:"foreignKey:OwnerId"`
}
