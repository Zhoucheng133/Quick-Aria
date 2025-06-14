package main

import (
	"database/sql"
	"quick_aria/utils"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2"
)

func main() {

	db, err := sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	utils.InitDB(db)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
