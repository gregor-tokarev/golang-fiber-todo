package router

import (
	"github.com/gofiber/fiber/v2"
	"goapi/controllers"
)

func InitUserRoutes(api fiber.Router) fiber.Router {
	group := api.Group("/auth")

	group.Post("/signup", controllers.Signup)
	group.Post("/login", controllers.Login)
	group.Post("/refresh", controllers.Refresh)

	return group
}
