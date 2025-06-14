package main

import (
	"quick_aria/utils"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2"
)

func Middleware(c *fiber.Ctx) error {

	path := c.Path()
	if strings.HasPrefix(path, "/api") {
		switch path {
		case "/api/noaccount", "/api/login", "/api/register":
			c.Next()
			return nil
		default:
			return utils.TokenCheck(c)
		}
	}

	err := c.Next()
	if err != nil {
		return err
	}

	return nil
}

func main() {

	utils.InitDB()
	utils.InitKey()

	app := fiber.New()

	app.Use(Middleware)
	app.Get("/api/noaccount", func(c *fiber.Ctx) error {
		return utils.NoAccount(c)
	})

	app.Listen(":3000")
}
