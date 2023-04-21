package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"goapi/models"
	"goapi/utils"
	"strconv"
)

func CreateTask(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)

	var task = &models.Task{}
	task.Status = "todo"
	task.Text = ""
	task.OwnerId = int(userId)

	DB.Create(&task)

	return ctx.Status(201).JSON(task)
}

func GetAllTasks(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)
	skip := ctx.Query("skip")
	take := ctx.Query("take")

	fmt.Println(skip, take)
	if len(skip) == 0 {
		skip = "0"
	}
	if len(take) == 0 {
		take = "10"
	}

	skipInt, err := strconv.Atoi(skip)
	takeInt, err := strconv.Atoi(take)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "wrong format",
		})
	}

	var tasks []models.Task
	DB.Preload("Owner").Where("owner_id = ?", int(userId)).Offset(skipInt).Limit(takeInt).Find(&tasks)

	return ctx.Status(200).JSON(tasks)
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
