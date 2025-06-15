package utils

import (
	"database/sql"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var Db *sql.DB
var JwtKey []byte

func InitKey() {
	// 初始化JWT密钥
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	JwtKey = []byte(id)
	// 测试代码，生产模式下注释下一行
	JwtKey = []byte("quick_aria")
}

func InitDB() {

	// 初始化数据库
	var err error
	Db, err = sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		panic(err)
	}
	// defer Db.Close()
	sqlInit(Db)
}

func sqlInit(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		name TEXT,
		password TEXT,
		aria TEXT,
		secret TEXT
    )
	`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}
