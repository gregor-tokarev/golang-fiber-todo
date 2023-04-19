package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"goapi/config"
	"goapi/models"
	"goapi/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB
var Validator = validator.New()

func init() {
	DB = models.InitDB()
}

func Signup(ctx *fiber.Ctx) error {
	var reqBody *models.SignupRequest
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "wrong format"})
	}

	err := Validator.Struct(reqBody)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": utils.CheckErrors(err)})
	}

	var user *models.User
	DB.Where("email = ?", reqBody.Email).First(&user)
	if user.Email != "" {
		return ctx.Status(400).JSON(fiber.Map{"message": "Email already exists"})
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(500).JSON(err)
	}

	user.Name = reqBody.Name
	user.Email = reqBody.Email
	user.Password = string(bytes)

	DB.Create(&user)

	t, err := generateToken(user.Id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(201).JSON(&models.SignupResponse{
		AccessToken: t,
	})
}

func Login(ctx *fiber.Ctx) error {
	var reqBody *models.LoginRequest
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "wrong format"})
	}

	err := Validator.Struct(reqBody)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": utils.CheckErrors(err)})
	}

	var user *models.User
	DB.Where("email = ?", reqBody.Email).First(&user)
	if user.Email == "" {
		return ctx.Status(400).JSON(fiber.Map{"message": "Email doesn't exist"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "Invalid password"})
	}

	t, err := generateToken(user.Id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(&models.LoginResponse{
		AccessToken: t,
	})
}

func generateToken(userId int) (string, error) {
	claims := jwt.MapClaims{}

	claims["sub"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Cfg.JwtAccessSecret))

	return tokenString, err
}
