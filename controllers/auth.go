package controllers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"goapi/config"
	"goapi/models"
	"goapi/models/provider"
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
	reqBody, err := utils.ValidateBody[models.SignupRequest](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": utils.CheckErrors(err)})
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

	t, err := generateTokens(user.Id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	err = updateRefreshToken(user.Id, t.RefreshToken)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(201).JSON(t)
}

func Login(ctx *fiber.Ctx) error {
	reqBody, err := utils.ValidateBody[models.LoginRequest](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": utils.CheckErrors(err)})
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

	t, err := generateTokens(user.Id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	err = updateRefreshToken(user.Id, t.RefreshToken)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(200).JSON(t)
}

func generateTokens(userId int) (models.Tokens, error) {
	accessClaims := jwt.MapClaims{}
	accessClaims["sub"] = userId
	accessClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	refreshClaims := jwt.MapClaims{}
	refreshClaims["sub"] = userId
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 14).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, accessErr := accessToken.SignedString([]byte(config.Cfg.JwtAccessSecret))
	refreshTokenString, refreshErr := refreshToken.SignedString([]byte(config.Cfg.JwtRefreshSecret))
	if accessErr != nil {
		return models.Tokens{}, accessErr
	}
	if refreshErr != nil {
		return models.Tokens{}, refreshErr
	}

	return models.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func Refresh(ctx *fiber.Ctx) error {
	userId, err := parseRefreshToken(ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "Invalid token"})
	}

	var user *models.User
	DB.Where("id = ?", userId).First(&user)
	if user.Email == "" {
		return ctx.Status(400).JSON(fiber.Map{"message": "Owner doesn't exist"})
	}

	body, err := utils.ValidateBody[models.RefreshRequest](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "error validating body"})
	}

	if user.RefreshToken != body.Token {
		return ctx.Status(401).JSON(fiber.Map{"message": "Invalid token"})
	}

	t, err := generateTokens(user.Id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	err = updateRefreshToken(user.Id, t.RefreshToken)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(200).JSON(t)
}

func updateRefreshToken(userId int, refreshToken string) error {
	var user *models.User
	DB.Where("id = ?", userId).First(&user)
	if user.Email == "" {
		return errors.New("Owner doesn't exist")
	}

	user.RefreshToken = refreshToken

	DB.Save(&user)

	return nil
}

func parseRefreshToken(ctx *fiber.Ctx) (int, error) {
	reqBody, err := utils.ValidateBody[models.RefreshRequest](ctx)
	if err != nil {
		return -1, err
	}

	token, err := jwt.Parse(reqBody.Token, jwtKeyFunc)
	if err != nil {
		return -1, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, err
	}
	if claims["sub"] == nil {
		return -1, err
	}

	return int(claims["sub"].(float64)), nil
}

func jwtKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(config.Cfg.JwtRefreshSecret), nil
}

func Logout(ctx *fiber.Ctx) error {
	userId := ctx.Locals("claims").(map[string]interface{})["sub"].(float64)

	var user *models.User
	DB.Where("id = ?", userId).First(&user)
	if user.Email == "" {
		return ctx.Status(400).JSON(fiber.Map{"message": "Owner doesn't exist"})
	}

	user.RefreshToken = ""
	DB.Save(&user)

	return ctx.Status(200).JSON(fiber.Map{"message": "Logged out"})
}

var googleProvider = provider.NewGoogleProvider(&provider.GoogleProviderConfig{

	ClientId:     config.Cfg.GoogleOauthClientID,
	ClientSecret: config.Cfg.GoogleOauthClientSecret,
	RedirectUri:  "http://localhost:3000/api/v1/auth/google/callback",
	Scope:        "profile+email",
})

func GoogleOauth(ctx *fiber.Ctx) error {
	uri, err := googleProvider.GetAuthUrl()
	fmt.Println(uri)
	if err != nil {
		panic(err)
	}

	return ctx.Redirect(uri)
}

func GoogleOauthCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")

	tokens, err := googleProvider.GetTokens(code)
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	userInfo, err := googleProvider.FetchInfo(tokens.AccessToken)
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	fmt.Println(userInfo)
	user := models.NewUserOauth(models.CreateUserOauthConfig{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Provider: "google",
	})

	authTokens, err := generateTokens(user.Id)
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Redirect(config.Cfg.FrontendUrl + fmt.Sprintf("?access_token=%s&refresh_token=%s", authTokens.AccessToken, authTokens.RefreshToken))
}
