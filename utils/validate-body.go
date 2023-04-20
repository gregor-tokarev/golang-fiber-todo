package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

func ValidateBody[T interface{}](ctx *fiber.Ctx) (*T, error) {
	var reqBody *T

	if err := ctx.BodyParser(&reqBody); err != nil {
		return nil, err
	}

	err := Validator.Struct(reqBody)
	if err != nil {
		return nil, err
	}

	return reqBody, nil
}
