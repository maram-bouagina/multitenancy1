package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// parseID extracts and validates a UUID from the route param ":id"
func ParseID(c *fiber.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	return id, err
}

// parseBody decodes the request body into dst and writes a 400 on failure
func ParseBody(c *fiber.Ctx, dst any) error {
	if err := c.BodyParser(dst); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	return nil
}

// fail writes a JSON error response
func Fail(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}
