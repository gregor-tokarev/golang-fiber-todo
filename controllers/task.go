package controllers

import (
	"github.com/gofiber/fiber/v2"
	"goapi/models"
	"goapi/utils"
	"gorm.io/gorm"
	"strconv"
)

func CreateTask(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)

	var task = &models.Task{}
	task.Status = "todo"
	task.Text = ""

	maxOrder := getMaxTaskOrder(int(userId))
	task.Order = maxOrder + 1

	task.OwnerId = int(userId)

	DB.Create(&task)

	return ctx.Status(201).JSON(task)
}

func GetAllTasks(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)
	skip := ctx.Query("skip")
	take := ctx.Query("take")

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
	reqBody, err := utils.ValidateBody[models.UpdateTask](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	taskId := ctx.Params("task_id")
	var task *models.Task
	DB.Where("id = ?", taskId).First(&task)

	task.Text = reqBody.Text
	task.DueDate = reqBody.DueDate
	task.Notes = reqBody.Notes

	DB.Save(&task)

	return ctx.Status(200).JSON(task)
}

func getMaxTaskOrder(userId int) int {
	var maxOrder int
	DB.Raw("SELECT max(tasks.order) FROM tasks WHERE owner_id = ?", userId).Scan(&maxOrder)
	return maxOrder
}

func ChangeTaskStatus(ctx *fiber.Ctx) error {
	reqBody, err := utils.ValidateBody[models.ChangeTaskStatus](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	taskId := ctx.Params("task_id")
	var task *models.Task

	DB.Where("id = ?", taskId).First(&task)
	task.Status = reqBody.Status
	if task.Status == "completed" {
		task.Order = -1
	}

	DB.Save(&task)

	return ctx.Status(200).JSON(task)
}

func ChangeTaskOrder(ctx *fiber.Ctx) error {
	reqBody, err := utils.ValidateBody[models.ChangeTaskOrder](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	taskId := ctx.Params("task_id")
	var task *models.Task
	DB.Where("id = ?", taskId).First(&task)

	initialOrder := task.Order

	if initialOrder > reqBody.Order {
		DB.
			Debug().
			Model(&models.Task{}).
			Where("\"owner_id\" = ? AND \"order\" >= ? AND \"order\" <= ?", task.OwnerId, reqBody.Order, initialOrder).
			Update("\"order\"", gorm.Expr("\"order\" + 1"))
	} else if initialOrder < reqBody.Order {
		DB.
			Debug().
			Model(&models.Task{}).
			Where("\"owner_id\" = ? AND \"order\" <= ? AND \"order\" > ?", task.OwnerId, reqBody.Order, initialOrder).
			Update("\"order\"", gorm.Expr("\"order\" - 1"))
	}

	task.Order = reqBody.Order
	DB.Save(&task)

	return ctx.Status(200).JSON(task)
}
