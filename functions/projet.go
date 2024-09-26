package functions

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func mod(a, b int) int {
	return a % b
}
func add(a, b int) int {
    return a + b
}

func Projet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var projectData []ProjectData
    rows, err := db.Query("SELECT title, description, link FROM projects")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    for rows.Next() {
        var project ProjectData
        err := rows.Scan(&project.Title, &project.Description, &project.Link)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
            return
        }
        projectData = append(projectData, project)
    }

    // Create a new template and add the custom function
    tmpl := template.New("projet.html").Funcs(template.FuncMap{
        "mod": mod,
		"add": add,
		
    })

    // Parse the template file
    tmpl, err = tmpl.ParseFiles("tmpl/projet.html")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
        return
    }

    // Execute the template with the project data
    err = tmpl.Execute(w, projectData)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
    }
}