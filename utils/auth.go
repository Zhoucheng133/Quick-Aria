package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

// 【GET】是否没有账户
func NoAccount(c *fiber.Ctx) error {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	return c.JSON(Response{
		Ok:  true,
		Msg: count == 0,
	})
}

// 【Func】验证密码
func checkPassword(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

// 【Func】生成JWT
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(365 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	return tokenString, err
}

// 【中间件】验证JWT
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
		return c.JSON(fiber.Map{
			"ok":  false,
			"msg": "无效或过期的 token",
		})
	}

	// token 有效，继续处理请求
	return c.Next()
}

// 【POST】登录
func Login(c *fiber.Ctx) error {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	if count == 0 {
		return c.JSON(Response{
			Ok:  false,
			Msg: "没有任何账户",
		})
	}

	var user LoginBody

	if err := c.BodyParser(&user); err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: "参数不正确",
		})
	}
	var password string

	err = Db.QueryRow("SELECT password FROM user WHERE name = ?", user.Name).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(Response{
				Ok:  false,
				Msg: "用户名不存在",
			})
		} else {
			return c.JSON(Response{
				Ok:  false,
				Msg: err.Error(),
			})
		}
	}
	if checkPassword(password, user.Password) {
		token, err := GenerateJWT(user.Name)
		if err != nil {
			return c.JSON(Response{
				Ok:  false,
				Msg: err.Error(),
			})
		}
		return c.JSON(Response{
			Ok:  true,
			Msg: token,
		})
	}

	return c.JSON(Response{
		Ok:  false,
		Msg: "密码不正确",
	})
}

// 【POST】注册
func Register(c *fiber.Ctx) error {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	if count != 0 {
		return c.JSON(Response{
			Ok:  false,
			Msg: "已有账户",
		})
	}

	var user RegisterBody

	if err := c.BodyParser(&user); err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: "参数不正确",
		})
	}
	stmt, err := Db.Prepare("INSERT INTO user (id, name, password, aria, secret) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	defer stmt.Close()
	id, err := gonanoid.New()
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	_, err = stmt.Exec(id, user.Name, string(hash), user.Aria, user.Secret)
	if err != nil {
		return c.JSON(Response{
			Ok:  false,
			Msg: err.Error(),
		})
	}
	return c.JSON(Response{
		Ok:  true,
		Msg: "",
	})
}
