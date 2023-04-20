package controllers

import (
	"github.com/gofiber/fiber/v2"
	"goapi/models"
	"goapi/utils"
)

func CreateTask(ctx *fiber.Ctx) error {
	var task = &models.Task{}
	task.Status = "todo"
	task.Text = ""

	DB.Create(&task)

	return ctx.Status(201).JSON(task)
}

func DeleteTask(ctx *fiber.Ctx) error {
	taskId := ctx.Params("task_id")

	DB.Delete(&models.Task{}, taskId)

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateTask(ctx *fiber.Ctx) error {
	var reqBody *models.UpdateTask
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "wrong format",
		})
	}

	err := Validator.Struct(reqBody)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	taskId := ctx.Params("task_id")
	var task *models.Task
	DB.Where("id = ?", taskId).First(&task)

	task.Text = reqBody.Text
	task.Status = reqBody.Status
	task.DueDate = reqBody.DueDate

	DB.Save(&task)

	return ctx.Status(200).JSON(task)
}
