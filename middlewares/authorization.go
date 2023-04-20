package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"goapi/config"
	"strings"
)

func AuthorizationAccess(ctx *fiber.Ctx) error {
	claims, err := ExtractDataFromAccessToken(ctx)
	if err != nil {
		return ctx.Status(403).JSON(fiber.Map{
			"message": "AuthorizationAccess failed",
		})
	}

	ctx.Locals("claims", claims)

	return ctx.Next()
}

func ExtractDataFromAccessToken(ctx *fiber.Ctx) (map[string]interface{}, error) {
	tokenString := extractAccessToken(ctx)
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

func extractAccessToken(ctx *fiber.Ctx) string {
	bearToken := ctx.Get("Authorization")

	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func jwtKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(config.Cfg.JwtAccessSecret), nil
}
