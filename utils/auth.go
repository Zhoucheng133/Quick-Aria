package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// 是否没有账户
func NoAccount(c *fiber.Ctx) error {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return c.JSON(fiber.Map{
			"ok":  false,
			"msg": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"ok":  true,
		"msg": count == 0,
	})
}

// 验证JWT
func TokenCheck(c *fiber.Ctx) error {
	tokenString := c.Get("token")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":  false,
			"msg": "缺少 token",
		})
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否符合预期
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":  false,
			"msg": "无效或过期的 token",
		})
	}

	// token 有效，继续处理请求
	return c.Next()
}
