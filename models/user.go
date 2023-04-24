package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id           int    `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	Provider     string `json:"provider"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken"`
	Tasks        []Task `gorm:"foreignKey:OwnerId"`
	Tags         []Tag  `gorm:"foreignKey:OwnerId"`
}

type CreateUserOauthConfig struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

func NewUserOauth(config CreateUserOauthConfig) *User {
	var user = &User{}
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)

	user.Provider = "local"

	DB.Create(&user)

	return user
}

func (u *User) Save(fields ...string) *User {
	DB.Select(fields).Save(u)
	return u
}

func FindUserByEmail(email string) *User {
	var user *User
	DB.Where("email = ?", email).First(&user)
	return user
}
func FindUserById(id int) *User {
	var user *User
	DB.Where("id = ?", id).First(&user)
	return user
}
