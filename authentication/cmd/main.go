package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"

	consulRegistry "github.com/content-management-system/backend/authentication/pkg/consul"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	client := consulRegistry.ConsulClient()
	consulRegistry.RegisterToConsul(client, consulRegistry.RegistrationService())
}
func main() {
	app := fiber.New(fiber.Config{
		AppName: "Go Fiber App",
	})
	log.Println("Hello World")
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://consul:8500,http://api-gateway:8080",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello Fiber!",
			"status":  "success",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
