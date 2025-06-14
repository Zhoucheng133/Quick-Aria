package utils

import "database/sql"

func InitDB(db *sql.DB) {
	initUser(db)
	initAria(db)
}

func initUser(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		username TEXT,
		password TEXT
    )
	`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func initAria(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS aria (
		id TEXT PRIMARY KEY,
		link TEXT,
		secret TEXT
    )
	`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}
