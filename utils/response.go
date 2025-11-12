package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

func RespondJSON(c *fiber.Ctx, statusCode int, status bool, message string, errors interface{}, data interface{}) error {
	response := Response{
		Status:  status,
		Message: message,
		Errors:  errors,
		Data:    data,
	}
	return c.Status(statusCode).JSON(response)
}

func FiberError(c *fiber.Ctx, msg string, err error) error {
	return RespondJSON(c, fiber.StatusInternalServerError, false, msg, []string{err.Error()}, nil)
}

func FiberSuccess(c *fiber.Ctx, msg string, data interface{}) error {
	return RespondJSON(c, fiber.StatusOK, true, msg, nil, data)
}

func FiberErrorCustom(c *fiber.Ctx, statusCode int, msg string, customErr string) error {
	return RespondJSON(c, statusCode, false, msg, []string{customErr}, nil)
}
