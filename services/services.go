package services

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	app.Post("/deploy", handleDeploy())
	app.Listen(":8080")
}
