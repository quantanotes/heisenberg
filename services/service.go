package main

import "github.com/gofiber/fiber/v2"

func RunService() {
	app := fiber.New()

	app.Get("/")
}
