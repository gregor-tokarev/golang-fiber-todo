package router

import (
	"github.com/gofiber/fiber/v2"
	"goapi/controllers"
	"goapi/middlewares"
)

func InitTagRouter(api fiber.Router) fiber.Router {
	group := api.Group("/tag")

	group.Use(middlewares.AuthorizationAccess)
	group.Get("/", controllers.GetAllTags)
	group.Post("/", controllers.CreateTag)
	group.Post("/:tag_id/connect/:task_id", controllers.ConnectTask)
	group.Post("/remove/:task_id", controllers.ClearTag)

	return group
}
