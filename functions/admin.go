package functions

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type UserData struct {
	ProfilePicture string
	Username       string
	Email          string
}

type ProjectData struct {
	Id          int
	Title       string
	Description string
	Link        string
	Picture     string
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
	var projectData []ProjectData
	rows, err := db.Query("SELECT id, title, description, link, picture FROM projects")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var project ProjectData
		var pictureByte []byte
		err := rows.Scan(&project.Id, &project.Title, &project.Description, &project.Link, &pictureByte)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
			return
		}
		if pictureByte != nil {
			project.Picture = base64.StdEncoding.EncodeToString(pictureByte)
		}
		projectData = append(projectData, project)
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

func CreateNewProject(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	//open the default profile picture
	file, err := os.Open("img/projectImg.jpg")
	if err != nil {
		return fmt.Errorf("error opening image: %v", err)
	}
	defer file.Close()

	//read the image
	PictureDefault, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading image: %v", err)
	}

	_, err = db.Exec("INSERT INTO projects (title, description, link, picture) VALUES (?, ?, ?, ?)", "", "", "", PictureDefault)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating a new project: %v", err), http.StatusInternalServerError)
		return nil
	}
	http.Redirect(w, r, "/editProjects", http.StatusSeeOther)
	return nil
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

	http.Redirect(w, r, "/editPersonal", http.StatusSeeOther)
}

func UpdateProjects(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusInternalServerError)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Error: no id provided", http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("picture")
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
		_, err = db.Exec("UPDATE projects SET picture=? WHERE id=?", fileBytes, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating the project picture: %v", err), http.StatusInternalServerError)
			return
		}
	}

	title := r.FormValue("title")
	if title != "" {
		_, err = db.Exec("UPDATE projects SET title=? WHERE id=?", title, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating the title: %v", err), http.StatusInternalServerError)
			return
		}
	}

	description := r.FormValue("description")
	if description != "" {
		_, err = db.Exec("UPDATE projects SET description=? WHERE id=?", description, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating the description: %v", err), http.StatusInternalServerError)
			return
		}
	}

	link := r.FormValue("link")
	if link != "" {
		_, err = db.Exec("UPDATE projects SET link=? WHERE id=?", link, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating the link: %v", err), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/editProjects", http.StatusSeeOther)
}
