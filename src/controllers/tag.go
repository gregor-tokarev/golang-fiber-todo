package controllers

import (
	"github.com/gofiber/fiber/v2"
	models2 "goapi/src/models"
	utils2 "goapi/src/utils"
)

func CreateTag(ctx *fiber.Ctx) error {
	reqBody, err := utils2.ValidateBody[models2.CreateTagReq](ctx)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"message": utils2.CheckErrors(err),
		})
	}

	userId := ctx.Locals("userId").(float64)

	tag := models2.NewTag(models2.NewTagConfig{
		Name:    reqBody.Name,
		OwnerId: int(userId),
	})

	return ctx.JSON(tag)
}

func ConnectTask(ctx *fiber.Ctx) error {
	tagId, err := ctx.ParamsInt("tag_id")
	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.JSON(fiber.Map{
			"message": "Invalid tag or task id",
		})
	}

	task := models2.FindTaskById(taskId)

	if task.Id == 0 {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Task doesn't exist",
		})
	}

	tag := models2.FindTagById(tagId)
	if tag.Id == 0 {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Tag doesn't exist",
		})
	}

	task.TagId = tagId

	task.Save("tag_id")

	return ctx.JSON(task)
}

func GetAllTags(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)

	tags := models2.FindAllTags(int(userId))

	return ctx.JSON(tags)
}

func ClearTag(ctx *fiber.Ctx) error {
	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.JSON(fiber.Map{
			"message": "Invalid task id",
		})
	}

	task := models2.FindTaskById(taskId)
	if task.Id == 0 {
		return ctx.JSON(fiber.Map{
			"message": "Task doesn't exist",
		})
	}

	task.ClearTag()

	return ctx.JSON(fiber.Map{
		"message": "Tag cleared",
	})
}
