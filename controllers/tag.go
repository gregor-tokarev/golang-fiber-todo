package controllers

import (
	"github.com/gofiber/fiber/v2"
	"goapi/models"
	"goapi/utils"
	"gorm.io/gorm"
)

func CreateTag(ctx *fiber.Ctx) error {
	reqBody, err := utils.ValidateBody[models.CreateTag](ctx)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	userId := ctx.Locals("userId").(float64)

	var tag models.Tag
	tag.Name = reqBody.Name
	tag.OwnerId = int(userId)
	DB.Create(&tag)

	return ctx.JSON(tag)
}

func ConnectTask(ctx *fiber.Ctx) error {
	tagId, err := ctx.ParamsInt("tag_id")
	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.JSON(fiber.Map{
			"message": "Invalid id",
		})
	}

	var task *models.Task
	DB.Where("id = ?", taskId).First(&task)

	if task.Id == 0 {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Task doesn't exist",
		})
	}

	var tag *models.Tag
	DB.Where("id = ?", tagId).First(&tag)
	if tag.Id == 0 {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Tag doesn't exist",
		})
	}

	task.TagId = tagId

	DB.Save(&task)

	return ctx.JSON(task)
}

func GetAllTags(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)

	var tags []models.Tag
	DB.Where("owner_id = ?", userId).Find(&tags)

	return ctx.JSON(tags)
}

func ClearTag(ctx *fiber.Ctx) error {
	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.JSON(fiber.Map{
			"message": "Invalid task id",
		})
	}

	var task models.Task
	DB.Where("id = ?", taskId).First(&task)
	if task.Id == 0 {
		return ctx.JSON(fiber.Map{
			"message": "Task doesn't exist",
		})
	}

	DB.Model(&task).Update("tag_id", gorm.Expr("NULL"))

	return ctx.JSON(fiber.Map{
		"message": "Tag cleared",
	})
}
