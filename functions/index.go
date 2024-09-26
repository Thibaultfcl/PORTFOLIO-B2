package functions

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./tmpl/index.html")
}

func Contact(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userData := UserData{}
	var userPPbyte []byte

	row := db.QueryRow("SELECT username, email, pp FROM user")

	err := row.Scan(&userData.Username, &userData.Email, &userPPbyte) // Exemple, il manque le contexte de `row`
	if err != nil {
		http.Error(w, fmt.Sprintf("Error scanning user data: %v", err), http.StatusInternalServerError)
		return
	}

	if userPPbyte != nil {
		userData.ProfilePicture = base64.StdEncoding.EncodeToString(userPPbyte)
	}

	tmpl, err := template.ParseFiles("./tmpl/contact.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Exécution du template avec les données utilisateur
	err = tmpl.Execute(w, userData)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
