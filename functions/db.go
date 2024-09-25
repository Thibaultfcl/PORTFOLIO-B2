package functions

import "database/sql"

func CreateTableUser(db *sql.DB) {
	//creating the user table if not already created
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username VARCHAR(12) NOT NULL,
            password VARCHAR(12) NOT NULL,
			email TEXT NOT NULL,
			isAdmin BOOL NOT NULL DEFAULT FALSE,
			isModerator BOOL NOT NULL DEFAULT FALSE,
			isBanned BOOL NOT NULL DEFAULT FALSE,
			pp BLOB,
			UUID VARCHAR(36) NOT NULL
        )
    `)
	if err != nil {
		panic(err.Error())
	}
}

func CreateTable(db *sql.DB) {
	CreateTableUser(db)
}