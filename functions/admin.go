package functions

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type UserData struct {
	ProfilePicture string
	Username       string
	Email          string
}

type ProjectData struct {
	Title       string
	Description string
	Link        string
}

func Admin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	http.ServeFile(w, r, "./tmpl/admin.html")
}

func EditPersonal(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userData := UserData{}
	row := db.QueryRow("SELECT username, email, pp FROM user")
	var userPPbyte []byte
	err := row.Scan(&userData.Username, &userData.Email, &userPPbyte)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	if userPPbyte != nil {
		userData.ProfilePicture = base64.StdEncoding.EncodeToString(userPPbyte)
	}

	tmpl, err := template.ParseFiles("tmpl/editPersonal.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, userData); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func EditProjects(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	projectData := ProjectData{}
	row := db.QueryRow("SELECT title, description, link FROM projects")
	err := row.Scan(&projectData.Title, &projectData.Description, &projectData.Link)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("tmpl/editProjects.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, projectData); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
