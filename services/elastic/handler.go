package elastic

import "github.com/gofiber/fiber/v2"

func CreateContainerHandler(c *fiber.Ctx) {
	usr := c.Params("usr")
	createContainer(usr)
}
