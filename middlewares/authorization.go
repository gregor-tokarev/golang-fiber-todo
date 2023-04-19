package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"goapi/config"
	"strings"
)

func Authorization(ctx *fiber.Ctx) error {
	claims, err := ExtractDataFromToken(ctx)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(403).JSON(fiber.Map{
			"message": "Authorization failed",
		})
	}

	fmt.Println(claims)

	return ctx.Next()
}

func ExtractDataFromToken(ctx *fiber.Ctx) (map[string]interface{}, error) {
	tokenString := extractToken(ctx)
	token, err := jwt.Parse(tokenString, jwtKeyFunc)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func jwtKeyFunc(_token *jwt.Token) (interface{}, error) {
	return []byte(config.Cfg.JwtAccessSecret), nil
}
