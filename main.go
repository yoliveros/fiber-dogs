package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

type Dog struct {
	Id    string
	Name  string
	Owner string
	Age   uint8
	Breed string
}

func handleDog(c *fiber.Ctx) error {
	dog := Dog{
		Id:    "test-id",
		Name:  "Kaiser",
		Owner: "Yoel",
		Age:   8,
		Breed: "Golden Retriever",
	}
	return c.Status(fiber.StatusOK).JSON(dog)
}

func handleCreateDog(c *fiber.Ctx) error {
	dog := Dog{}
	if err := c.BodyParser(&dog); err != nil {
		return err
	}

	dog.Id = uuid.NewString()

	return c.Status(fiber.StatusCreated).JSON(dog)
}

func main() {
	app := fiber.New()

	// Middleware
	app.Use(logger.New(), requestid.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	dogGroup := app.Group("/dogs")
	dogGroup.Get("", handleDog)
	dogGroup.Post("", handleCreateDog)

	app.Listen(":3000")
}
