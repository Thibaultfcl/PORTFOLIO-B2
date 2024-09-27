package functions

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Projet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var projectsData []ProjectData
	rows, err := db.Query("SELECT title, description, link, picture FROM projects")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var project ProjectData
		var pictureByte []byte
		err := rows.Scan(&project.Title, &project.Description, &project.Link, &pictureByte)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
			return
		}

		if pictureByte != nil {
			project.Picture = base64.StdEncoding.EncodeToString(pictureByte)
		}

		projectsData = append(projectsData, project)
	}

	tmpl, err := template.ParseFiles("tmpl/projet.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, projectsData); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
