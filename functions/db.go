package functions

import (
	"database/sql"
	"fmt"
	"io"
	"os"
)

func CreateTableUser(db *sql.DB) {
	//creating the user table if not already created
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS user (
            username VARCHAR(12) NOT NULL,
			email TEXT NOT NULL,
			pp BLOB
        )
    `)
	if err != nil {
		panic(err.Error())
	}
}

func CreateTableProjects(db *sql.DB) {
	//creating the projects table if not already created
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			link TEXT NOT NULL
		)
	`)
	if err != nil {
		panic(err.Error())
	}
}

func SetUserDefault(db *sql.DB) error {
	//check if the user account already exists
	row := db.QueryRow("SELECT COUNT(*) FROM user")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("error while checking the user account: %v", err)
	}
	if count > 0 {
		return nil
	}

	//open the default profile picture
	file, err := os.Open("img/profileDefault.jpg")
	if err != nil {
		return fmt.Errorf("error opening image: %v", err)
	}
	defer file.Close()

	//read the image
	ppDefault, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading image: %v", err)
	}

	//insert the admin account in the database
	_, err = db.Exec("INSERT INTO user (username, email, pp) VALUES (?, ?, ?)", "user", "user@email.com", ppDefault)
	if err != nil {
		return fmt.Errorf("error while creating the user account: %v", err)
	}

	return nil
}

func SetProjectsDefault(db *sql.DB) error {
	//check if the projects already exist
	row := db.QueryRow("SELECT COUNT(*) FROM projects")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("error while checking the projects: %v", err)
	}
	if count > 0 {
		return nil
	}

	//insert the default projects in the database
	_, err = db.Exec("INSERT INTO projects (title, description, link) VALUES (?, ?, ?)", "project", "description", "link")
	if err != nil {
		return fmt.Errorf("error while creating the projects: %v", err)
	}

	return nil
}

func CreateTable(db *sql.DB) {
	CreateTableUser(db)
	CreateTableProjects(db)
}
