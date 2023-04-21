package router

import (
	"github.com/gofiber/fiber/v2"
	"goapi/controllers"
	"goapi/middlewares"
)

func InitTaskRouter(api fiber.Router) fiber.Router {
	group := api.Group("/task")

	group.Use(middlewares.AuthorizationAccess)

	// Crate new task
	group.Post("/", controllers.CreateTask)
	// update task details
	group.Put("/:task_id", controllers.UpdateTask)
	// delete task
	group.Delete("/:task_id", controllers.DeleteTask)
	group.Get("/", controllers.GetAllTasks)

	return group
}
