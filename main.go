package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dog struct {
	gorm.Model
	Id    string
	Name  string
	Owner string
	Age   uint8
	Breed string
}

var Database *gorm.DB
var DB_NAME = "dogs.db"

func handleDog(c *fiber.Ctx) error {
	var dogs []Dog

	Database.Find(&dogs)

	return c.Status(fiber.StatusOK).JSON(dogs)
}

func handleCreateDog(c *fiber.Ctx) error {
	dog := Dog{}
	if err := c.BodyParser(&dog); err != nil {
		return err
	}

	dog.Id = uuid.NewString()

	Database.Create(&dog)

	return c.Status(fiber.StatusCreated).JSON(dog)
}

func main() {
	// DB Connection
	db, err := gorm.Open(sqlite.Open("dogs.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&Dog{})

	app := fiber.New()

	// Middleware
	app.Use(logger.New(), requestid.New(), cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	dogGroup := app.Group("/dogs")
	dogGroup.Get("", handleDog)
	dogGroup.Post("", handleCreateDog)

	app.Listen(":3000")
}
