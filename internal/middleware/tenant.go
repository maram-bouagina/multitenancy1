package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"multitenancypfe/internal/database"
)

func TenantDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID, ok := c.Locals("userID").(string)
		if !ok || tenantID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "tenant not identified"})
		}

		schema := fmt.Sprintf("tenant_%s", tenantID)

		sqlDB, err := database.DB.DB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
		}

		scopedDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db scope error"})
		}

		scopedDB = scopedDB.Exec(fmt.Sprintf("SET search_path TO %q, public", schema))

		c.Locals("tenantDB", scopedDB)
		c.Locals("tenantID", tenantID)
		c.Locals("schema", schema)

		return c.Next()
	}
}

func GetTenantDB(c *fiber.Ctx) *gorm.DB {
	db, _ := c.Locals("tenantDB").(*gorm.DB)
	return db
}
