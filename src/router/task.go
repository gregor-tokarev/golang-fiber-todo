package router

import (
	"github.com/gofiber/fiber/v2"
	"goapi/src/controllers"
	"goapi/src/middlewares"
)

func InitTaskRouter(api fiber.Router) fiber.Router {
	group := api.Group("/task")

	group.Use(middlewares.AuthorizationAccess)

	// Crate new task
	group.Post("/", controllers.CreateTask)
	// update task details
	group.Put("/:task_id", controllers.UpdateTask)
	group.Patch("/:task_id/status", controllers.ChangeTaskStatus)
	group.Patch("/:task_id/order", controllers.ChangeTaskOrder)
	// delete task
	group.Delete("/:task_id", controllers.DeleteTask)
	group.Get("/", controllers.GetAllTasks)

	return group
}
