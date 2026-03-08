package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"multitenancypfe/internal/config"
	"multitenancypfe/internal/database"
	"multitenancypfe/internal/jwt"
	"multitenancypfe/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config failed: %v", err)
	}

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("database failed: %v", err)
	}

	jwt.Init(cfg.JWTSecret)

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	routes.Register(app, database.DB)

	log.Printf("API running on :%s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
