package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"portfolio/functions"

	_ "github.com/mattn/go-sqlite3"
)

// port of the server
const port = ":8080"

func main() {
	//open the database with sqlite3
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err.Error())
	}

	//create the tables
	functions.CreateTable(db)
	functions.SetUserDefault(db)
	functions.SetProjectsDefault(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { functions.Index(w, r) })
	http.HandleFunc("/projet", func(w http.ResponseWriter, r *http.Request) { functions.Projet(w, r) })
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) { functions.Contact(w, r) })
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) { functions.Admin(w, r, db) })
	http.HandleFunc("/editProjects", func(w http.ResponseWriter, r *http.Request) { functions.EditProjects(w, r, db) })
	http.HandleFunc("/editPersonal", func(w http.ResponseWriter, r *http.Request) { functions.EditPersonal(w, r, db) })
	http.HandleFunc("/updateUserInfo", func(w http.ResponseWriter, r *http.Request) { functions.UpdateUserInfo(w, r, db) })

	//load the CSS, the JS and the images
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./script"))))

	//start the local host
	fmt.Println("\n(http://localhost:8080/) - Server started on port", port)
	http.ListenAndServe(port, nil)
}
