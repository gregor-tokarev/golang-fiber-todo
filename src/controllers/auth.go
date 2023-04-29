package controllers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"goapi/src/config"
	models2 "goapi/src/models"
	"goapi/src/models/provider"
	utils2 "goapi/src/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Signup(ctx *fiber.Ctx) error {
	reqBody, err := utils2.ValidateBody[models2.SignupReq](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": utils2.CheckErrors(err)})
	}

	user := models2.NewUser(models2.NewUserConfig{
		Email:    reqBody.Email,
		Name:     reqBody.Name,
		Password: reqBody.Password,
	})
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
	reqBody, err := utils2.ValidateBody[models2.LoginReq](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": utils2.CheckErrors(err)})
	}

	user := models2.FindUserByEmail(reqBody.Email)
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

func generateTokens(userId int) (models2.Tokens, error) {
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
		return models2.Tokens{}, accessErr
	}
	if refreshErr != nil {
		return models2.Tokens{}, refreshErr
	}

	return models2.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func Refresh(ctx *fiber.Ctx) error {
	userId, err := parseRefreshToken(ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"message": "Invalid token"})
	}

	user := models2.FindUserById(userId)
	if user.Email == "" {
		return ctx.Status(400).JSON(fiber.Map{"message": "Owner doesn't exist"})
	}

	body, err := utils2.ValidateBody[models2.RefreshReq](ctx)
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
	user := models2.FindUserById(userId)
	if user.Email == "" {
		return errors.New("Owner doesn't exist")
	}

	user.RefreshToken = refreshToken

	user.Save("refresh_token")

	return nil
}

func parseRefreshToken(ctx *fiber.Ctx) (int, error) {
	reqBody, err := utils2.ValidateBody[models2.RefreshReq](ctx)
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

	user := models2.FindUserById(int(userId))
	if user.Email == "" {
		return ctx.Status(400).JSON(fiber.Map{"message": "Owner doesn't exist"})
	}

	user.RefreshToken = ""
	user.Save("refresh_token")

	return ctx.Status(200).JSON(fiber.Map{"message": "Logged out"})
}

var googleProvider = provider.NewGoogleProvider(&provider.GoogleProviderConfig{
	ClientId:     config.Cfg.GoogleOauthClientID,
	ClientSecret: config.Cfg.GoogleOauthClientSecret,
	RedirectUri:  config.Cfg.ServerHost + config.Cfg.ServerPrefix + "/auth/google/callback",
	Scope:        "profile+email",
})

func GoogleOauth(ctx *fiber.Ctx) error {
	uri, err := googleProvider.GetAuthUrl()
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

	user := models2.NewUserOauth(models2.CreateUserOauthConfig{
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

	err = updateRefreshToken(user.Id, authTokens.RefreshToken)
	if err != nil {
		return ctx.Status(500).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Redirect(config.Cfg.FrontendUrl + fmt.Sprintf("?access_token=%s&refresh_token=%s", authTokens.AccessToken, authTokens.RefreshToken))
}
