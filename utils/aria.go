package utils

import "github.com/gofiber/fiber/v2"

func GetAria(c *fiber.Ctx) error {
	var data AriaData

	err := Db.QueryRow("SELECT aria, secret FROM user").Scan(&data.Aria, &data.Secret)
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	return c.JSON(Response{
		Ok: true,
		Msg: fiber.Map{
			"aria":   data.Aria,
			"secret": data.Secret,
		},
	})
}
