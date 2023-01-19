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
	Id    string `gorm:"primarykey"`
	Name  string
	Owner string
	Age   uint8
	Breed string
}

var db *gorm.DB
var DB_NAME = "dogs.db"

func getDog(c *fiber.Ctx) error {
	var dogs []Dog

	db.Find(&dogs)

	return c.Status(fiber.StatusOK).JSON(dogs)
}

func addDog(c *fiber.Ctx) error {
	dog := Dog{}
	if err := c.BodyParser(&dog); err != nil {
		return err
	}

	dog.Id = uuid.NewString()

	db.Create(&dog)

	return c.Status(fiber.StatusCreated).JSON(dog)
}

func updateDog(c *fiber.Ctx) error {
	dog := Dog{}
	id := c.Params("id")

	if err := c.BodyParser(&dog); err != nil {
		return err
	}

	db.Where("id = ?", id).Updates(&dog)
	return c.Status(fiber.StatusOK).JSON(dog)
}

func deleteDog(c *fiber.Ctx) error {
	id := c.Params("id")

	result := db.Delete(&Dog{}, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.SendStatus(fiber.StatusOK)
}

func main() {
	// DB Connection
	var err error
	db, err = gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
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
	dogGroup.Get("", getDog)
	dogGroup.Post("", addDog)
	dogGroup.Patch("", updateDog)
	dogGroup.Delete("", deleteDog)

	app.Listen(":3000")
}
