package controllers

import (
	"github.com/gofiber/fiber/v2"
	"goapi/models"
	"goapi/utils"
	"strconv"
)

func CreateTask(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)

	task := models.NewTask(models.CreateTaskConfig{
		OwnerId: int(userId),
	})

	return ctx.Status(201).JSON(task)
}

func GetAllTasks(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(float64)
	skip := ctx.Query("skip")
	take := ctx.Query("take")
	tagId := ctx.Query("tag_id")

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

	tasks := models.FindTasks(models.FindTasksConfig{
		OwnerId: int(userId),
		Skip:    skipInt,
		Take:    takeInt,
		TagId:   tagId,
	})

	return ctx.Status(200).JSON(tasks)
}

func DeleteTask(ctx *fiber.Ctx) error {
	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "wrong format",
		})
	}

	task := models.FindTaskById(taskId)
	task.Delete()

	return ctx.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateTask(ctx *fiber.Ctx) error {
	reqBody, err := utils.ValidateBody[models.UpdateTaskReq](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "wrong task_id format",
		})
	}
	task := models.FindTaskById(taskId)

	task.Text = reqBody.Text
	task.DueDate = reqBody.DueDate
	task.Notes = reqBody.Notes

	task.Save("text", "due_date", "notes")

	return ctx.Status(200).JSON(task)
}

func ChangeTaskStatus(ctx *fiber.Ctx) error {
	reqBody, err := utils.ValidateBody[models.ChangeTaskStatusReq](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "wrong format",
		})
	}
	task := models.FindTaskById(taskId)

	task.Status = reqBody.Status
	if task.Status == "completed" {
		task.Order = -1
	} else if task.Status == "todo" {
		task.Order = 1
	}

	task.Save("status", "order")

	return ctx.Status(200).JSON(task)
}

func ChangeTaskOrder(ctx *fiber.Ctx) error {
	reqBody, err := utils.ValidateBody[models.ChangeTaskOrderReq](ctx)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": utils.CheckErrors(err),
		})
	}

	taskId, err := ctx.ParamsInt("task_id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "wrong format",
		})
	}
	task := models.FindTaskById(taskId)

	task = task.ChangeTaskOrder(reqBody.Order)

	return ctx.Status(200).JSON(task)
}
