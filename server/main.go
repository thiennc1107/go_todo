package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title,omitempty"`
	Done  bool   `json:"done"`
	Body  string `json:"body,omitempty"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("../todo.db"), &gorm.Config{})

	getAll := func() ([]Todo, error) {
		var todos []Todo
		result := db.Find(&todos)
		if result.Error != nil {
			return nil, result.Error
		}
		return todos, nil
	}

	db.AutoMigrate(&Todo{})

	if err != nil {
		panic("Failed to connect db:" + err.Error())
	}

	fmt.Println("Initializing server")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPatch,
		}, ","),
	}))

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(401).SendString("Invalid id")
		}
		var todo Todo
		result := db.Model(&todo).Where("ID = ?", id).Update("Done", true)
		if result.Error != nil {
			return result.Error
		}

		todos, err := getAll()
		if err != nil {
			return err
		}

		return c.JSON(todos)
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		todos, err := getAll()
		if err != nil {
			return err
		}
		return c.JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}
		result := db.Create(&todo)

		if result.Error != nil {
			return result.Error
		}
		todos, err := getAll()
		if err != nil {
			return err
		}
		return c.JSON(todos)
	})

	log.Fatal(app.Listen(":4000"))
}
