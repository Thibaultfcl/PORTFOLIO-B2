package functions

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
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

func UpdateUserInfo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("profilePicture")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, fmt.Sprintf("Error retrieving the file: %v", err), http.StatusInternalServerError)
			return
		}
	}
	if file != nil {
		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading the file: %v", err), http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE user SET pp=?", fileBytes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating the profile picture: %v", err), http.StatusInternalServerError)
			return
		}
	}

	username := r.FormValue("username")
	if username != "" {
		_, err = db.Exec("UPDATE user SET username=?", username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating the username: %v", err), http.StatusInternalServerError)
			return
		}
	}

	email := r.FormValue("email")
	if email != "" {
		_, err = db.Exec("UPDATE user SET email=?", email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating the email: %v", err), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
