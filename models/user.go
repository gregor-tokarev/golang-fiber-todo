package models

type User struct {
	Id           int    `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	Provider     string `json:"provider"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken"`
	Tasks        []Task `gorm:"foreignKey:OwnerId"`
}

type CreateUserOauthConfig struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

func NewUserOauth(config CreateUserOauthConfig) *User {
	var user *User
	DB.Where("email = ?", config.Email).First(&user)
	if user.Email != "" {
		return user
	}

	user.Email = config.Email
	user.Name = config.Name
	user.Provider = config.Provider

	DB.Create(&user)

	return user
}

type NewUserConfig struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(config NewUserConfig) *User {
	var user *User
	DB.Where("email = ?", config.Email).First(&user)
	if user.Email != "" {
		return user
	}

	user.Name = config.Name
	user.Email = config.Email
	user.Password = config.Password
	user.Provider = "local"

	DB.Create(&user)

	return user
}

func (u *User) Save() *User {
	DB.Save(&u)
	return u
}
